package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"tkOptimizer/internal/evolution"
	"tkOptimizer/internal/keyboard"
)

func main() {
	var (
		ymlPath string
	)
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

  // initial setup of 0 generation of keyboards
	if len(e.KeyboardConfig.Layout) != 0 {
		fmt.Println("initial layout:")
    kFromLayout := keyboard.NewEmpty(
      e.KeyboardConfig.Height,
      e.KeyboardConfig.Width,
      keyboard.SetLayout(e.KeyboardConfig.Layout),
      keyboard.SetWeights(e.KeyboardConfig.Weights),
    )
    kFromLayout.Print()
		loopKeyboards = []*keyboard.Keyboard{kFromLayout}
	} else {
		loopKeyboards, err = evolution.GenerateKeyboardsThreads(
      e.Threads, 
      e.KeyboardConfig,
      e.GetInitPopulation(),
    )
		if err != nil {
			panic(err)
		}
	}
	evolution.TestKeyboardsThreads(e.Threads, loopKeyboards, e.TestText)
	evolution.SortKeyboards(loopKeyboards)
	e.AppendMetric(loopKeyboards[0].Distance)
	bestK = loopKeyboards[0]

  accumulationK := make([]*keyboard.Keyboard, 0)
	i := 0
	loop, ok := true, true
  // main loop of finding optiomal keyboard
	for loop {
		fmt.Printf(
			"\rGen: %d, Last metric: %.2f",
			i,
			e.GetMetricLast(),
		)

		loopKeyboards, err = evolution.RecombineWithOne(
			loopKeyboards,
      e.MutationProbability,
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
		}

		if e.MinPopulation != 0 && len(loopKeyboards) > e.MinPopulation {
			loopKeyboards, ok = evolution.FilterPopulation(
        loopKeyboards, 
        e.Percentile, 
        e.MinPopulation,
      )
		}

		if !ok {
      if e.MinPopulation != 0 && e.MinPopulation == 1 {
        loop = false 
      }
			if len(accumulationK) != 0 {
				if loopKeyboards[0].Distance < accumulationK[0].Distance {
					accumulationK = append(accumulationK, loopKeyboards[0])
				}
			} else {
				accumulationK = append(accumulationK, bestK)
			}
			if len(accumulationK) > 20 {

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
        )
				if err != nil {
					panic(err)
				}
			}
		}

		i++

		select {
		case <-sigs:
			// If Ctrl+C is pressed, exit the loop
			fmt.Println("\nReceived Ctrl+C. Exiting...")
			loop = false
		default:
		}
	}
  // closing and cleaning up
	close(sigs)

	fmt.Println("Best distance: ", bestK.Distance)
	bestK.Print()
	bestK.PrintYamlFormat()
}
