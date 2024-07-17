package keyboard

import (
  "github.com/goccy/go-yaml"
  "os"
)

type PreConfigYaml struct {
  CharSet string `yaml:"characters"`
  Weights [][]float64 `yaml:"weights"`
  Layout [][]string `yaml:"layout"`
}

type ConfigYaml struct {
  CharSet []rune
  Weights [][]float64
  // Layout TODO transformation from Pre to Config
}

func FromYaml(filePath string) (*ConfigYaml, error) {
  // error on not existance 
  _, err := os.Stat(filePath)
  if err != nil {
    return nil, err
  }

  file, err := os.Open(filePath)
  if err != nil {
    return nil, err
  }

  var c PreConfigYaml
  err = yaml.NewDecoder(file).Decode(&c)
  if err != nil {
    return nil, err
  }

  return &ConfigYaml{
    CharSet: []rune(c.CharSet),
    Weights: c.Weights,
  }, nil
}
