package evolution

import (
	"tkOptimizer/internal/keyboard"
)

func (e *Evolution) RecombineWithOne(k []*keyboard.Keyboard, one *keyboard.Keyboard) ([]*keyboard.Keyboard, error) {
  variantsNew := make([]*keyboard.Keyboard, len(k))
  for i, k2 := range k {
    kM, err := Recombination(e.MutationProbability, one, k2)
    if err != nil {
      return nil, err
    }
    variantsNew[i] = kM
  }
  return variantsNew, nil
}

func (e *Evolution) RecombineWithOneThreads(threads int, k []*keyboard.Keyboard, one *keyboard.Keyboard) ([]*keyboard.Keyboard, error) {
  variantsNew := make([]*keyboard.Keyboard, len(k))
  errorChan := make(chan error, 10)
  jobs := make(chan *keyboard.Keyboard, len(k))
  results := make(chan *keyboard.Keyboard, GetLenOr1000(len(k)))

  for w := 1; w <= threads; w++ {
    go func() {
      for j := range jobs {
        kM, err := Recombination(e.MutationProbability, one, j)
        if err != nil {
          errorChan <- err
          return
        }
        results <- kM
      }
    }()
  }

  for i := 0; i < len(k); i++ {
    jobs <- k[i]
  }

  for j := 0; j < len(k); j++ {
    select {
    case err := <-errorChan:
      return nil, err
    case kN := <-results:
      variantsNew[j] = kN
    }
  }

  close(errorChan)
  close(jobs)
  close(results)

  return variantsNew, nil
}

func (e *Evolution) FilterPopulation(k []*keyboard.Keyboard, percentile float64) ([]*keyboard.Keyboard, bool) {
  filteredNumber := PopulationSizeNext(percentile, len(k))
  if !IsEven(filteredNumber) {
    filteredNumber--
  }
  if filteredNumber == 0 {
    return k, false
  }
  filtered := make([]*keyboard.Keyboard, filteredNumber)
  for i := 0; i < filteredNumber; i++ {
    filtered[i] = k[i]
  }
  
  clear(k)
  return filtered, true
}
