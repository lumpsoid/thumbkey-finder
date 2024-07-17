package evolution

import (
	"fmt"
	"math"
	"tkOptimizer/internal/keyboard"
)

func (e *Evolution) GenerateKeyboards(population int) error {
	height := e.KeyboardConfig.Height
	width := e.KeyboardConfig.Width
	weights := e.KeyboardConfig.Weights
	charSet := e.KeyboardConfig.CharSet

  keyboardsNew := make([]*keyboard.Keyboard, 0)
	for i := 0; i < population; i++ {
		var k *keyboard.Keyboard
		var err error

		if len(weights) != 0 {
			k, err = keyboard.GenerateNewWithWeights(
				height, width, weights, charSet,
			)
			if err != nil {
				return err
			}
		} else {
			k, err = keyboard.GenerateNew(height, width, charSet)
			if err != nil {
				return err
			}
		}
    keyboardsNew = append(keyboardsNew, k)
	}

  e.SetKeyboards(keyboardsNew)
	return nil
}

func (e *Evolution) TestKeyboards() {
	for i := 0; i < len(e.Keyboards); i++ {
    e.Keyboards[i].TravelDistance(e.TestText)
	}
  if len(e.Keyboards) > 1 {
    e.SortKeyboards(e.Keyboards)
  }
	e.AppendMetric(e.Keyboards[0].Distance)
}

func (e *Evolution) Recombine() error {
	passNumber := math.Floor(e.Percentile * float64(len(e.Keyboards)))
	passNumberInt := int(passNumber)
	if !IsEven(passNumberInt) {
		passNumberInt -= 1
	}
	fmt.Println("passNumber of keyboards: ", passNumberInt)
	if passNumberInt < 2 {
		return nil
	}

	nextGen := make([]*keyboard.Keyboard, 0)
	fmt.Println("NextGen keyboards count before: ", len(nextGen))
	for i := 0; i < passNumberInt; i += 2 {
		mK, err := Recombination(
			e.MutationProbability,
			e.Keyboards[i],
			e.Keyboards[i+1],
		)
		if err != nil {
			return err
		}
		nextGen = append(nextGen, mK)
	}
	e.Keyboards = nextGen
	fmt.Println("NextGen keyboards count after: ", len(nextGen))

	return nil
}

func (e *Evolution) Run() error {
  err := e.GenerateKeyboards(e.GetInitPopulation())
  if err != nil {
    return err
  }
	for len(e.Keyboards) > 2 {
		fmt.Println("Gen: ", len(e.MetricHistory))
		e.TestKeyboards()
		fmt.Println("Best distance: ", e.MetricHistory[len(e.MetricHistory)-1])
		err = e.Recombine()
    if err != nil {
      return err
    }
		fmt.Println("Keyboards count: ", len(e.Keyboards))
	}
  return nil
}
