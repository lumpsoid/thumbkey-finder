package keyboard

import (
	"fmt"
	"testing"
	"tkOptimizer/internal/key"
	"tkOptimizer/internal/layout"
)

func TestRandomCharInsert(t *testing.T) {
	characters := []rune("abcdefgklmnoprst")
	k := NewEmpty(9,9)
  k.GenerateWeigths(1.0)
	err := k.RandomCharInsertSafe(characters)
	//k.Print()
	if err != nil {
		t.Error(err)
	}
	return
}

func TestTravelDistance(t *testing.T) {
	k := NewEmpty(
    8,8,
		SetLayout(layout.Parse(map[string]key.Position{
      "r": {X: 1, Y: 8},
      "a": {X: 8, Y: 1},
    })),
	)
  k.GenerateWeigths(1.0)

  rightDistance := 15.462185
  k.TravelDistance("ra")
  if fmt.Sprintf("%.4f", k.Distance) != fmt.Sprintf("%.4f", rightDistance ) {
    t.Errorf("Expected %4f, got %4f distance", rightDistance, k.Distance)
  }
}
