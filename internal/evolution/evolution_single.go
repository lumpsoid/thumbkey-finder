package evolution

import (
	"tkOptimizer/internal/keyboard"
)

func GenerateKeyboards(
  config *KeyboardConfiguration, 
  population int,
  placeThreshold float64,
) ([]*keyboard.Keyboard, error) {
	height := config.Height
	width := config.Width
	weights := config.Weights
	charSet := config.CharSet

  genNew := make([]*keyboard.Keyboard, population)
	for i := 0; i < population; i++ {
		var k *keyboard.Keyboard
		var err error

		if len(weights) != 0 {
			k, err = keyboard.GenerateNewWithWeights(
				height, width, weights, charSet, placeThreshold,
			)
			if err != nil {
				return nil, err
			}
		} else {
			k, err = keyboard.GenerateNew(height, width, charSet, placeThreshold)
			if err != nil {
				return nil, err
			}
		}
    genNew[i] = k
	}

	return genNew, nil
}

func TestKeyboards(k []*keyboard.Keyboard, testText string) {
	for i := 0; i < len(k); i++ {
    k[i].TravelDistance(testText)
	}
}

func Recombine(
  k []*keyboard.Keyboard, 
  mutationProbability float64,
  placeThreshold float64,
) ([]*keyboard.Keyboard, error) {
  keyboardLen := len(k)
	if !IsEven(keyboardLen) {
		keyboardLen -= 1
	}

	nextGen := make([]*keyboard.Keyboard, 0)
	for i := 0; i < keyboardLen; i += 2 {
		mK, err := Recombination(
			mutationProbability,
      placeThreshold,
			k[i],
			k[i+1],
		)
		if err != nil {
			return nil, err
		}
		nextGen = append(nextGen, mK)
	}
	return nextGen, nil
}

func Run(e *Evolution, k []*keyboard.Keyboard) ([]*keyboard.Keyboard, error) {
	var err error
	var ok bool

  bestK := k[0]

	for len(k) > 2 {
    TestKeyboards(k, e.TestText)
    SortKeyboards(k)
		e.AppendDistance(k[0].Distance)

		if e.Threads == 1 {
			k, err = RecombineWithOne(k, e.MutationProbability, e.PlaceThreshold, bestK)
		} else {
			k, err = RecombineWithOneThreads(e.Threads, k, e.MutationProbability, e.PlaceThreshold, bestK)
		}

		if err != nil {
			return nil, err
		}

		if len(k) > e.MinPopulation {
      if e.Threads == 1 {
        k, err = Recombine(k, e.MutationProbability, e.PlaceThreshold)
      } else {
        k, err = RecombineThreads(e.Threads, e.MutationProbability, e.PlaceThreshold,k)
      }
		}

		if !ok {
			if e.MinPopulation == 1 {
				return k, nil
			}

			if e.Threads == 1 {
				k, err = GenerateKeyboards(e.KeyboardConfig, e.GetInitPopulation(), e.PlaceThreshold)
			} else {
				k, err = GenerateKeyboardsThreads(e.Threads, e.KeyboardConfig, e.GetInitPopulation(), e.PlaceThreshold)
			}

			if err != nil {
				return nil, err
			}
		}
	}

	return k, nil
}
