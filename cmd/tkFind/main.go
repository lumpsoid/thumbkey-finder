package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tkOptimizer/internal/evolution"
	"tkOptimizer/internal/keyboard"
	"tkOptimizer/internal/weights"
)

func main() {
	var ymlPath string
	flag.StringVar(&ymlPath, "config", "", "path to a yaml config for run")
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	e, err := evolution.FromYaml(ymlPath)
	if err != nil {
		panic(err)
	}

	var bestK *keyboard.Keyboard
	var loopKeyboards []*keyboard.Keyboard

	fmt.Println("Generating first generation...")
	loopKeyboards, err = evolution.GenerateKeyboardsThreads(
		e.Threads,
		e.KeyboardConfig,
		e.GetInitPopulation(),
    e.PlaceThreshold,
	)
	if err != nil {
		panic(err)
	}

	// initial setup of 0 generation of keyboards
	if len(e.KeyboardConfig.Layout) != 0 {
		fmt.Println("initial layout:")
		kFromLayout := keyboard.NewEmpty(
			e.KeyboardConfig.Height,
			e.KeyboardConfig.Width,
			keyboard.SetLayout(e.KeyboardConfig.Layout),
			keyboard.SetWeights(e.KeyboardConfig.Weights),
		)
		if len(kFromLayout.Weights) == 0 {
			kFromLayout.Weights = weights.New(
				kFromLayout.GetHeight(),
				kFromLayout.GetWidth(), 1,
			)
		}
		kFromLayout.Print()
		kFromLayout.TravelDistance(e.TestText)
		e.AppendMetric(kFromLayout.Distance)
		bestK = kFromLayout
	} else {
		bestK = loopKeyboards[0]
	}

	evolution.TestKeyboardsThreads(e.Threads, loopKeyboards, e.TestText)
	evolution.SortKeyboards(loopKeyboards)
	e.AppendMetric(loopKeyboards[0].Distance)

	var endTime time.Time
	startTime := time.Now()
	accumulationK := make([]*keyboard.Keyboard, 0)
	i, staleCounter := 0, 0
	loop, ok := true, true
	// main loop of finding optiomal keyboard
	for loop {
		fmt.Printf(
			//"\rGen: %d, Last metric: %.2f",
			"\rGen: %d, Len: %d, Metric: %.2f     ",
			i,
			len(loopKeyboards),
			e.GetMetricLast(),
		)

		loopKeyboards, err = evolution.RecombineWithOne(
			loopKeyboards,
			e.MutationProbability,
      e.PlaceThreshold,
			bestK,
		)
		if err != nil {
			panic(err)
		}
		evolution.TestKeyboardsThreads(e.Threads, loopKeyboards, e.TestText)
		evolution.SortKeyboards(loopKeyboards)
		e.AppendMetric(loopKeyboards[0].Distance)

		if loopKeyboards[0].Distance < bestK.Distance {
			bestK = loopKeyboards[0]
			fmt.Printf("\nDistance: %.2f\n", bestK.Distance)
			bestK.Print()
		} else {
      staleCounter++
    }

		if len(loopKeyboards) > e.MinPopulation {
			loopKeyboards, ok = evolution.FilterPopulation(
				loopKeyboards,
				e.Percentile,
				e.MinPopulation,
			)
		}

		if !ok || staleCounter > e.StaleThreshold {
			if e.MinPopulation == 1 {
				loop = false
			}
			if len(accumulationK) != 0 {
				if loopKeyboards[0].Distance*1.2 < accumulationK[0].Distance {
					accumulationK = append(accumulationK, loopKeyboards[0])
				}
			} else {
				accumulationK = append(accumulationK, bestK)
			}
			if len(accumulationK) > 10 {

				transferK := make([]*keyboard.Keyboard, len(accumulationK))
				for j, tK := range accumulationK {
					transferK[j] = tK
				}
				copy(transferK, loopKeyboards)
				accumulationK = make([]*keyboard.Keyboard, 0)
				clear(transferK)
			} else {
				loopKeyboards, err = evolution.GenerateKeyboardsThreads(
					e.Threads,
					e.KeyboardConfig,
					e.GetInitPopulation(),
          e.PlaceThreshold,
				)
				if err != nil {
					panic(err)
				}
			}
      staleCounter = 0
		}

		i++

		select {
		case <-sigs:
			endTime = time.Now()
			// If Ctrl+C is pressed, exit the loop
			fmt.Println("\nReceived Ctrl+C. Exiting...")
			loop = false
			fmt.Printf("Loop time: %.4f min\n", endTime.Sub(startTime).Minutes())
		default:
		}
	}
	// closing and cleaning up
	close(sigs)

	fmt.Println("Best distance: ", bestK.Distance)
	bestK.Print()
	fmt.Println("Export keyboard: ")
	bestK.PrintYamlFormat()
}
