package evolution

import (
	"tkOptimizer/internal/keyboard"
)

func RecombineWithOne(
  k []*keyboard.Keyboard, 
  mutationProbability float64, 
  placeThreshold float64, 
  one *keyboard.Keyboard,
) ([]*keyboard.Keyboard, error) {
  variantsNew := make([]*keyboard.Keyboard, len(k))
  for i, k2 := range k {
    kM, err := Recombination(mutationProbability, placeThreshold, one, k2)
    if err != nil {
      return nil, err
    }
    variantsNew[i] = kM
  }
  return variantsNew, nil
}

func RecombineWithOneThreads(
  threads int, 
  k []*keyboard.Keyboard, 
  mutationProbability float64, 
  placeThreshold float64, 
  one *keyboard.Keyboard,
) ([]*keyboard.Keyboard, error) {
  variantsNew := make([]*keyboard.Keyboard, len(k))
  errorChan := make(chan error, 10)
  jobs := make(chan *keyboard.Keyboard, len(k))
  results := make(chan *keyboard.Keyboard, GetLenOr1000(len(k)))

  for w := 1; w <= threads; w++ {
    go func() {
      for j := range jobs {
        kM, err := Recombination(mutationProbability, placeThreshold, one, j)
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

func FilterPopulation(
  k []*keyboard.Keyboard, 
  percentile float64, 
  minPopulation int,
) ([]*keyboard.Keyboard, bool) {
  filteredNumber := PopulationSizeNext(percentile, len(k))
  if filteredNumber <= 1 {
    return k, false
  }
  if !IsEven(filteredNumber) {
    filteredNumber--
  }
  filtered := make([]*keyboard.Keyboard, filteredNumber)
  for i := 0; i < filteredNumber; i++ {
    filtered[i] = k[i]
  }
  
  clear(k)
  return filtered, true
}
