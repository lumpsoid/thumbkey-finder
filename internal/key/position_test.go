package key

import (
	"testing"
)

func TestIsCentralPosition(t *testing.T) {
  posCheck := Position{X: 1, Y: 1}
  isCentral := isCentralPosition(posCheck, mapOfCentralPositions)
  if isCentral != true {
    t.Error("Exptected to be true")
    return
  }
  posCheck = Position{X: 8, Y: 7}
  isCentral = isCentralPosition(posCheck, mapOfCentralPositions)
  if isCentral != false {
    t.Error("Exptected to be false")
    return
  }
}


func TestClosestPosition(t *testing.T) {
  posCheck := Position{X: 1, Y: 1}
	startPosition := closestPosition(posCheck, mapOfCentralPositions)
  if startPosition != posCheck {
    t.Error(startPosition, posCheck)
    return
  }
  posCheck = Position{X: 1, Y: 2}
	startPosition = closestPosition(posCheck, mapOfCentralPositions)
  rightPosition := Position{X: 1, Y: 1}
  if startPosition != rightPosition {
    t.Error(startPosition, posCheck)
    return
  }
  posCheck = Position{X: 3, Y: 3}
	startPosition = closestPosition(posCheck, mapOfCentralPositions)
  rightPosition = Position{X: 4, Y: 4}
  if startPosition != rightPosition {
    t.Error(startPosition, posCheck)
    return
  }
}

func TestIsBelowDiagonal(t *testing.T) {
	positions := []Position{
		{X: 1, Y: 1},
		{X: 4, Y: 4},
		{X: 4, Y: 3},
		{X: 4, Y: 7},
		{X: 7, Y: 2},
	}
	rightAnswers := []bool{
    false,
		true,
		false,
		true,
		true,
	}
	for i := range positions {
		k := New(positions[i])
		rightAnswer := rightAnswers[i]

		below := k.Position.IsBelowDiagonal(8)
		if below != rightAnswer {
			t.Errorf(
				"Expected %t, got %t on position %s",
				rightAnswer,
				below,
				k.Position.String(),
			)
		}
	}
}
