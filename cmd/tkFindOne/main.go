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
  err = e.RunThreads()
  if err != nil {
    panic(err)
  }
  e.TestKeyboards()
  e.Keyboards[0].Print()
  fmt.Println("Best distance: ", e.Keyboards[0].Distance)
	e.Keyboards[0].Print()
	e.Keyboards[0].PrintYamlFormat()
}
