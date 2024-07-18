package main

import (
	"flag"
	"fmt"
	"tkOptimizer/internal/evolution"
	"tkOptimizer/internal/keyboard"
)

func main() {
  var (
    ymlPath string
  )
  flag.StringVar(&ymlPath, "config", "", "path to a yaml config for run")
  flag.Parse()

  e, err := evolution.FromYaml(ymlPath)
  if err != nil {
    panic(err)
  }
  k := keyboard.NewEmpty(
    e.KeyboardConfig.Height,
    e.KeyboardConfig.Width,
    keyboard.SetLayout(e.KeyboardConfig.Layout),
    keyboard.SetWeights(e.KeyboardConfig.Weights),
  )
  k.TravelDistance(e.TestText)
  k.Print()
  fmt.Println("Distanse: ", k.Distance)
}
