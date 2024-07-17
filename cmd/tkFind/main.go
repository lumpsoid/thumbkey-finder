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

  if len(e.Keyboards) != 0 {
    fmt.Println("init layout get")
    e.Keyboards[0].Print()
    loopKeyboards = e.Keyboards
  } else {
    loopKeyboards, err = e.GenerateKeyboardsThreads(e.GetInitPopulation())
    if err != nil {
      panic(err)
    }
    e.SetKeyboards(loopKeyboards)
  }
  e.TestKeyboardsThreads(loopKeyboards)
  e.SortKeyboards(loopKeyboards)
  e.AppendMetric(loopKeyboards[0].Distance)
  // TODO KeyboardConfig.Layout start search from giver layout
  bestK = loopKeyboards[0]

	accumulationK := make([]*keyboard.Keyboard, 0)
	i := 0
	loop, ok := true, true
	for loop {
		fmt.Printf(
			"\rGen: %d, Last metric: %.2f ''",
			i,
			e.GetMetricLast(),
		)

		loopKeyboards, err = e.RecombineWithOne(
			loopKeyboards,
			bestK,
		)
		if err != nil {
			panic(err)
		}
		e.TestKeyboardsThreads(loopKeyboards)
		e.SortKeyboards(loopKeyboards)
		e.AppendMetric(loopKeyboards[0].Distance)

		if loopKeyboards[0].Distance < bestK.Distance {
			bestK = loopKeyboards[0]
			fmt.Printf("\nDistance: %.2f\n", bestK.Distance)
			bestK.Print()
		}

		loopKeyboards, ok = e.FilterPopulation(loopKeyboards, e.Percentile)

		e.SetKeyboards(loopKeyboards)

		if !ok {
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
				loopKeyboards, err = e.GenerateKeyboardsThreads(e.GetInitPopulation())
				if err != nil {
					panic(err)
				}
				e.SetKeyboards(loopKeyboards)
				e.SetPopulation(e.GetInitPopulation())
			}
		}

		i++

		select {
		case <-sigs:
			// If Ctrl+C is pressed, exit the loop
			fmt.Println("\nReceived Ctrl+C, exiting loop")
			loop = false
		default:
		}
	}

	bestK.Print()
	fmt.Println("Best distance: ", bestK.Distance)

	// Handle cleanup or final tasks here if needed
	fmt.Println("Exiting...")
}
