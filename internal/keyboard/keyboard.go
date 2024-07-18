package keyboard

import (
	"errors"
	"fmt"
	"math/rand"
	"tkOptimizer/internal/key"
)

type Weights [][]float64

func (w *Weights) Check(height int, width int) error {
	rows := *w
	if len(rows) == 0 {
		return nil
	}
	if len(rows)-1 != height {
		return errors.New("weights rows are not equal to keyboard height")
	}
	if len(rows[0])-1 != width {
		return errors.New("weights cols are not equal to keyboard height")
	}
	return nil
}

type Keyboard struct {
	height    int
	width     int
	Layout    map[rune]*key.Key
	Weights   Weights
	Positions map[int]map[int]int
	Distance  float64
}

// cmp(a, b) should return a negative number when a \< b, a positive number when a > b and zero when a == b
func SortCMP(a *Keyboard, b *Keyboard) int {
	if a.Distance < b.Distance {
		return -1
	}
	if a.Distance > b.Distance {
		return 1
	}
	return 0
}

func SortCMPDes(a *Keyboard, b *Keyboard) int {
	if a.Distance < b.Distance {
		return 1
	}
	if a.Distance > b.Distance {
		return -1
	}
	return 0
}

func (k *Keyboard) GetHeight() int {
	return k.height
}

func (k *Keyboard) GetWidth() int {
	return k.width
}

func (k *Keyboard) GetKeyByIndex(index int) (rune, *key.Key) {
	i := 0
	// Iterate over the map to find the key at the specified index
	for c, k := range k.Layout {
		if i == index {
			return c, k
		}
		i++
	}
	return 0, nil
}

func (k *Keyboard) GetCharByPosition(pos key.Position) (rune, *key.Key) {
	// Iterate over the map to find the key at the specified index
	for c, k := range k.Layout {
		if pos == k.Position {
			return c, k
		}
	}
	return 0, nil
}

func (k *Keyboard) GetRandomKey() (rune, *key.Key) {
	// Intn [0, len(k.Layout))
	randomIndex := rand.Intn(len(k.Layout))
	return k.GetKeyByIndex(randomIndex)
}

func (k *Keyboard) SwapChars(
	c1 rune, k1 *key.Key,
	c2 rune, k2 *key.Key,
) {
	k.Layout[c1] = k2
	k.Layout[c2] = k1
	return
}

func (k *Keyboard) GenerateWeigths(filler float64) {
	weightsNew := make([][]float64, k.height)

	for i := 0; i < len(weightsNew); i++ {
		weightsNew[i] = make([]float64, k.width)
		for j := 0; j < len(weightsNew[0]); j++ {
			weightsNew[i][j] = filler
		}
	}
	k.Weights = weightsNew
	return
}

func NewEmpty(height int, width int, options ...KeyboardOption) *Keyboard {
	// Create a slice of slices with 9 inner slices
	layout := make(map[rune]*key.Key)
	weights := make([][]float64, height)
	positions := make(map[int]map[int]int)

	keyboard := &Keyboard{
		height:    height,
		width:     width,
		Layout:    layout,
		Weights:   weights,
		Positions: positions,
	}
	applyOptions(keyboard, options...)
	return keyboard
}

func probabilityToPlace(weight float64) bool {
	return rand.Float64()*weight > 0.5
}

func (k *Keyboard) PositionExist(pos key.Position) bool {
	_, rowExists := k.Positions[int(pos.Y)]
	if !rowExists {
		return false
	}
	_, cellExists := k.Positions[int(pos.Y)][int(pos.X)]
	return cellExists
}

func (k *Keyboard) InsertPosition(pos key.Position) error {
	if k.PositionExist(pos) {
		return AlreadyExist{Message: "position already exists"}
	}

	rowIndex := int(pos.Y)
	cellIndex := int(pos.X)

	if _, ok := k.Positions[rowIndex]; !ok {
		k.Positions[rowIndex] = make(map[int]int)
	}

	k.Positions[rowIndex][cellIndex] = 1
	return nil
}

func (k *Keyboard) InsertKey(charNew rune, keyNew *key.Key) error {
	err := k.InsertPosition(keyNew.Position)
	if err != nil {
		fmt.Println("from insertKey pos: ", keyNew.Position, " char: ", string(charNew))
		fmt.Println(int(keyNew.Position.X))
		return err
	}
	k.Layout[charNew] = keyNew
	return nil
}

func (k *Keyboard) RandomCharInsertSafe(charSlice []rune) error {
	rows := len(k.Weights)
	if rows == 0 {
		return errors.New("weights are empty")
	}
	cols := len(k.Weights[0])
	if cols == 0 {
		return errors.New("cols in weights are empty")
	}
	r, c := 0, -1

	for i := 0; i < len(charSlice); {
		c++
		if c == cols {
			c = 0
			r++
			if r == rows {
				r = 0
			}
		}
		char := charSlice[i]
		placeIt := probabilityToPlace(k.Weights[r][c])

		if placeIt {
			positionNew := key.Position{X: float64(c), Y: float64(r)}

			if k.PositionExist(positionNew) {
				// Position exists, skip this iteration
				continue
			}

			// Attempt to insert the character
			err := k.InsertKey(char, key.New(positionNew))
			if err != nil {
				if IsAlreadyExist(err) {
					continue
				}
				return err
			}

			// Move to the next character in charSlice
			i++
		}
	}

	return nil
}

func ShuffleSlice(data []rune) []rune {
	// Create a copy of the input slice
	shuffledData := make([]rune, len(data))
	copy(shuffledData, data)
	// Fisher-Yates shuffle algorithm
	for i := len(shuffledData) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		shuffledData[i], shuffledData[j] = shuffledData[j], shuffledData[i]
	}
	return shuffledData
}

func GenerateNew(height int, width int, charSet []rune) (*Keyboard, error) {
	k := NewEmpty(height, width)
	k.GenerateWeigths(1.0)

	characters := ShuffleSlice(charSet)
	err := k.RandomCharInsertSafe(characters)
	if err != nil {
		return nil, err
	}

	return k, nil
}

func GenerateNewWithWeights(
	height int,
	width int,
	weights Weights,
	charSet []rune,
) (*Keyboard, error) {
	k := NewEmpty(height, width)
	k.Update(SetWeights(weights))

	characters := ShuffleSlice(charSet)
	err := k.RandomCharInsertSafe(characters)
	if err != nil {
		return nil, err
	}

	return k, nil
}

func GenerateFromYaml(c *ConfigYaml) (*Keyboard, error) {
	layout := make(map[rune]*key.Key)
	positions := make(map[int]map[int]int)

	k := &Keyboard{
		Layout:    layout,
		Weights:   c.Weights,
		Positions: positions,
	}

	characters := ShuffleSlice(c.CharSet)
	err := k.RandomCharInsertSafe(characters)
	if err != nil {
		return nil, err
	}

	return k, nil
}

func (k *Keyboard) GetKeyByChar(c rune) (*key.Key, error) {
	v, ok := k.Layout[c]
	if !ok {
		return nil, errors.New("error char is not in the layout")
	}
	return v, nil
}

func (k *Keyboard) TravelDistance(text string) {
	distance := 0.0
	// X0.0 Y0.0 is top left
	prevPos := key.Position{X: 4.0, Y: 4.0} // Center Starting position

	for _, char := range text {
		currKey, err := k.GetKeyByChar(char)
		if err != nil {
			// Ignore characters not in the layout
			continue
		}

		if currKey.Type == key.Press {
			distance += key.ComputeDistance(prevPos, currKey.Position)
		} else if currKey.Type == key.Swipe {
			distance += key.ComputeDistance(prevPos, currKey.Central)
			distance += key.ComputeDistance(currKey.Central, currKey.Position)
		} else {
			panic(fmt.Errorf("error: Key type is not Press or Swipe"))
		}

		prevPos = currKey.Position
	}

	k.Distance = distance
	return
}

func printDashLine(length int) {
	for i := 0; i < length; i++ {
		fmt.Print("-")
	}
}

func (k *Keyboard) Print() {
	layout := make([][]rune, 9)

	for i := 0; i < len(layout); i++ {
		layout[i] = make([]rune, 9)
		for j := 0; j < len(layout[0]); j++ {
			layout[i][j] = '·'
		}
	}

	for char, key := range k.Layout {
		layout[int(key.Position.Y)][int(key.Position.X)] = char
	}

	for row := range layout {
		if row != 0 && row%3 == 0 {
			printDashLine(9 + 2)
			fmt.Println()
		}
		for char := range layout[row] {
			if char != 0 && char%3 == 0 {
				fmt.Print("|")
			}
			fmt.Printf("%v", string(layout[row][char]))
		}
		fmt.Println()
	}
}

func (k *Keyboard) PrintYamlFormat() {
	layout := make([][]string, 9)

	for i := 0; i < len(layout); i++ {
		layout[i] = make([]string, 9)
		for j := 0; j < len(layout[0]); j++ {
			layout[i][j] = ""
		}
	}

	for char, key := range k.Layout {
		layout[int(key.Position.Y)][int(key.Position.X)] = string(char)
	}

  fmt.Println("[")
	for row := range layout {
    fmt.Print("  [")
		for char := range layout[row] {
			fmt.Printf(`"%s",`, layout[row][char])
		}
		fmt.Println("],")
	}
  fmt.Println("]")
}
