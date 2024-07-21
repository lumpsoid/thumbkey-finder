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
	var generetionCurrent []*keyboard.Keyboard

	fmt.Println("Generating first generation...")
	generetionCurrent, err = evolution.GenerateKeyboardsThreads(
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
		e.AppendDistance(kFromLayout.Distance)
		bestK = kFromLayout
	} else {
		bestK = generetionCurrent[0]
	}

	evolution.TestKeyboardsThreads(e.Threads, generetionCurrent, e.TestText)
	evolution.SortKeyboards(generetionCurrent)
	e.AppendDistance(generetionCurrent[0].Distance)

	var endTime time.Time
	startTime := time.Now()
	memoryPool := make([]*keyboard.Keyboard, 0)
	i, staleCounter := 0, 0
	loop, ok := true, true
	// main loop of finding optiomal keyboard
	for loop {
		fmt.Printf(
			//"\rGen: %d, Last metric: %.2f",
			"\rGen: %d, Len: %d, Metric: %.2f     ",
			i,
			len(generetionCurrent),
			e.GetMetricLast(),
		)

		generetionCurrent, err = evolution.RecombineWithOne(
			generetionCurrent,
			e.MutationProbability,
			e.PlaceThreshold,
			bestK,
		)
		if err != nil {
			panic(err)
		}
		evolution.TestKeyboardsThreads(e.Threads, generetionCurrent, e.TestText)
		evolution.SortKeyboards(generetionCurrent)
		e.AppendDistance(generetionCurrent[0].Distance)

		if generetionCurrent[0].Distance < bestK.Distance {
			bestK = generetionCurrent[0]
			fmt.Printf("\nDistance: %.2f\n", bestK.Distance)
			bestK.Print()
		} else {
			staleCounter++
		}

		if len(generetionCurrent) > e.MinPopulation {
			generetionCurrent, ok = evolution.FilterPopulationSafe(
				1,
				generetionCurrent,
				e.Percentile,
				e.MinPopulation,
			)
		}

		if !ok || staleCounter > e.StaleThreshold {
			if e.MinPopulation == 1 {
				loop = false
			}
			if len(memoryPool) != 0 {
				if generetionCurrent[0].Distance*1.2 < memoryPool[0].Distance {
					memoryPool = append(memoryPool, generetionCurrent[0])
				}
			} else {
				memoryPool = append(memoryPool, bestK)
			}
			if len(memoryPool) > 10 {

				transferK := make([]*keyboard.Keyboard, len(memoryPool))
				for j, tK := range memoryPool {
					transferK[j] = tK
				}
				copy(transferK, generetionCurrent)
				memoryPool = make([]*keyboard.Keyboard, 0)
				clear(transferK)
			} else {
				generetionCurrent, err = evolution.GenerateKeyboardsThreads(
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
