package evolution

import (
	"errors"
	"math"
	"math/rand"
	"slices"
	"tkOptimizer/internal/key"
	"tkOptimizer/internal/keyboard"
)

// Height and Width are 0 based X and Y coordinates
type KeyboardConfiguration struct {
	Height  int
	Width   int
	Weights keyboard.Weights
	CharSet []rune
}

// height and width are 0 based X and Y coordinates
// for 9x9 keyboard it should be 8 and 8
func NewKeyboardConfig(
	height int,
	width int,
	weights [][]float64,
	charSet []rune,
) *KeyboardConfiguration {
	return &KeyboardConfiguration{
		Height:  height,
		Width:   width,
		Weights: weights,
		CharSet: charSet,
	}
}

type Evolution struct {
	Threads             int
	initPopulation      int
	Population          int
	Percentile          float64
  MinPopulation       int
	MutationProbability float64
	KeyboardConfig      *KeyboardConfiguration
	Keyboards           []*keyboard.Keyboard // in the descending order
	TestText            string
	MetricHistory       []float64
}

func New(
	threads int,
	numKeyboards int,
	persentile float64,
	mutationProbability float64,
	config *KeyboardConfiguration,
	textString string,
) (*Evolution, error) {
	if numKeyboards%2 != 0 {
		return nil, errors.New("even number of init keyboards")
	}
	history := make([]float64, 0)
	return &Evolution{
		Threads:             threads,
		initPopulation:      numKeyboards,
		MutationProbability: mutationProbability,
		Population:          0,
		Percentile:          persentile,
		KeyboardConfig:      config,
		TestText:            textString,
		MetricHistory:       history,
	}, nil
}

func mergeKeyboards(
	keyboard1 *keyboard.Keyboard,
	keyboard2 *keyboard.Keyboard,
) (*keyboard.Keyboard, error) {
	keyboardMerged := keyboard.NewEmpty(
		keyboard1.GetHeight(),
		keyboard1.GetWidth(),
		keyboard.SetWeights(keyboard1.Weights),
	)
	kHeight := keyboard1.GetHeight()

	needInsert := make([]rune, 0)
	for char, key := range keyboard1.Layout {
		keyFrom2, err := keyboard2.GetKeyByChar(char)
		if err != nil {
			return nil, err
		}

		// taking char if it is from bottom part of 1d keyboard
		if key.Position.IsBelowDiagonal(kHeight) {
			err := keyboardMerged.InsertKey(char, key)
			if err != nil {
				return nil, err
			}
			// taking char if it is from upper part of 2d keyboard
		} else if !keyFrom2.Position.IsBelowDiagonal(kHeight) {
			err := keyboardMerged.InsertKey(char, keyFrom2)
			if err != nil {
				return nil, err
			}
			// if char form opposite parts in both variants,
			// save for futher processing
		} else {
			needInsert = append(needInsert, char)
		}
	}
	// insert chars which are lefted after merge both sides
	needInsert = keyboard.ShuffleSlice(needInsert)
	err := keyboardMerged.RandomCharInsertSafe(needInsert)
	if err != nil {
		return nil, err
	}
	return keyboardMerged, nil
}

func probabilityToSwap(probability float64) bool {
	return rand.Float64() < probability
}

func Mutation(mergedKeyboard *keyboard.Keyboard, probability float64) {
	for c1, k1 := range mergedKeyboard.Layout {
		if probabilityToSwap(probability) {
			var c2 rune
			var k2 *key.Key
			for k2 == nil {
				c2, k2 = mergedKeyboard.GetRandomKey()
			}

			mergedKeyboard.SwapChars(c1, k1, c2, k2)
		}
	}
}

// keyboard1 more performant variant
func Recombination(
	probability float64,
	keyboard1 *keyboard.Keyboard,
	keyboard2 *keyboard.Keyboard,
) (*keyboard.Keyboard, error) {
	if keyboard1.GetHeight() != keyboard2.GetHeight() {
		return nil, errors.New("not equal height in keyboard1 and keyboard2")
	}
	keyboardMerged, err := mergeKeyboards(keyboard1, keyboard2)
	if err != nil {
		return nil, err
	}

	Mutation(keyboardMerged, probability)

	return keyboardMerged, nil
}

func (e *Evolution) SortKeyboards(k []*keyboard.Keyboard) {
	slices.SortFunc(k, keyboard.SortCMPDes)
}

func (e *Evolution) SetKeyboards(k []*keyboard.Keyboard) {
	e.Keyboards = k
}

func (e *Evolution) AppendKeyboard(k *keyboard.Keyboard) {
	e.Keyboards = append(e.Keyboards, k)
}

func (e *Evolution) GetInitPopulation() int {
	return e.initPopulation
}

func PopulationSizeNext(percentile float64, population int) int {
	return int(math.Ceil(float64(population) * percentile))
}

func (e *Evolution) SetPopulation(num int) {
	e.Population = num
}

func (e *Evolution) AppendMetric(result float64) {
	e.MetricHistory = append(e.MetricHistory, result)
}

func (e *Evolution) GetMetricLast() float64 {
	return e.MetricHistory[len(e.MetricHistory)-1]
}

func IsEven(num int) bool {
	return num%2 == 0
}
