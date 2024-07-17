package main

import (
	"flag"
	"fmt"
	"tkOptimizer/internal/evolution"
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
  e.TestKeyboards()
  e.Keyboards[0].Print()
  fmt.Println("Distanse: ", e.Keyboards[0].Distance)
}
