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
	var generationCurrent []*keyboard.Keyboard

	fmt.Println("Generating first generation...")
	generationCurrent, err = evolution.GenerateKeyboardsThreads(
		e.Threads,
		e.KeyboardConfig,
		e.GetInitPopulation(),
		e.PlaceThreshold,
	)
	if err != nil {
		panic(err)
	}

	// if we have prefered layout from yaml
  // set up to use it as first best keyboard
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
		bestK = generationCurrent[0]
	}

	evolution.TestKeyboardsThreads(e.Threads, generationCurrent, e.TestText)
	evolution.SortKeyboards(generationCurrent)
	e.AppendDistance(generationCurrent[0].Distance)

	var endTime time.Time
	startTime := time.Now()
	memoryPool := make([]*keyboard.Keyboard, 0)
	i, staleCounter, lastBestGeneration := 0, 0, 0
	loop := true
  // loop to optimize layout
	for loop {
		fmt.Printf(
			//"\rGen: %d, Last metric: %.2f",
			"\rGen: %d, Len: %d, Metric: %.2f     ",
			i,
			len(generationCurrent),
			e.GetMetricLast(),
		)
    // recombine with current generation using best keyboard
    // as performant variant
    // output: same population count
    // but with best keyboard information incorporated
		generationCurrent, err = evolution.RecombineWithOneThreads(
			e.Threads,
			generationCurrent,
			e.MutationProbability,
			e.PlaceThreshold,
			bestK,
		)
		if err != nil {
			panic(err)
		}

		evolution.TestKeyboardsThreads(e.Threads, generationCurrent, e.TestText)
		evolution.SortKeyboards(generationCurrent)
		e.AppendDistance(generationCurrent[0].Distance)

    // if population count greater than specified minimum population in yaml
    // recombine them between themselves 
    // 2 parents : 1 child
		if len(generationCurrent) > e.MinPopulation {
			generationCurrent, err = evolution.RecombineThreads(
				e.Threads,
				e.MutationProbability,
				e.PlaceThreshold,
				generationCurrent,
			)
			evolution.TestKeyboardsThreads(e.Threads, generationCurrent, e.TestText)
			evolution.SortKeyboards(generationCurrent)
			e.AppendDistance(generationCurrent[0].Distance)
		}

		if generationCurrent[0].Distance < bestK.Distance {
			bestK = generationCurrent[0]
			fmt.Printf("\nDistance: %.2f\n", bestK.Distance)
			fmt.Printf("Last best update (generations): %d\n", lastBestGeneration - i)
			bestK.Print()
      lastBestGeneration = i
      
		} else {
			staleCounter++
		}

    // reset bestK keyboard if local optimization is great
    if i - lastBestGeneration > e.ResetThreshold  {
			fmt.Print("\nReseting best keyboard because of the stale\n")
      bestK = generationCurrent[0]
      lastBestGeneration = i
			fmt.Printf("New distance: %.2f\n", bestK.Distance)
      bestK.Print()

    }

		if staleCounter > e.StaleThreshold {
			if len(memoryPool) == 0 {
				memoryPool = append(memoryPool, bestK)
			}
			// if keyboard's distance in current generation
			// in 20% range from best
			// it will saved in memory pool

			// TODO check all current generation until >
			// not only first
			if generationCurrent[0].Distance < memoryPool[0].Distance*1.2 {
				memoryPool = append(memoryPool, generationCurrent[0])
			}
			generationCurrent, err = evolution.GenerateKeyboardsThreads(
				e.Threads,
				e.KeyboardConfig,
				e.GetInitPopulation(),
				e.PlaceThreshold,
			)
			if err != nil {
				panic(err)
			}
			// use memory pool to incorporate information from performant variants
      // to random generated ones
      // TODO make weights back propagation like technique
			var memoryPoolBestK *keyboard.Keyboard

			if len(memoryPool) >= 2 {
				memoryPoolBest, err := evolution.RecombineThreads(
					e.Threads,
					e.MutationProbability,
					e.PlaceThreshold,
					memoryPool,
				)
				if err != nil {
					panic(err)
				}
				evolution.TestKeyboardsThreads(e.Threads, memoryPoolBest, e.TestText)
				evolution.SortKeyboards(memoryPoolBest)
			  memoryPool = memoryPoolBest
			}
			memoryPoolBestK = memoryPool[0]
			generationCurrent, err = evolution.RecombineWithOneThreads(
				e.Threads,
				generationCurrent,
				e.MutationProbability,
				e.PlaceThreshold,
				memoryPoolBestK,
			)
			if err != nil {
				panic(err)
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
