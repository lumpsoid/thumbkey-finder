package key

import (
	"fmt"
	"math"
)

type Position struct {
	X, Y float64
}

func (p *Position) String() string {
  return fmt.Sprint("X: ", p.X, " Y: ", p.Y)
}

func (p *Position) IsBelowDiagonal(height int) bool {
  // Adjust the calculation to check if Y is greater than or equal to the diagonal line
	return p.Y >= float64(height) - p.X
}

var mapOfCentralPositions = map[Position]bool{
	{1, 1}: true,
	{4, 1}: true,
	{7, 1}: true,
	{1, 4}: true,
	{4, 4}: true,
	{7, 4}: true,
	{1, 7}: true,
	{4, 7}: true,
	{7, 7}: true,
}

func closestPosition(current Position, positions map[Position]bool) Position {
	var closestPos Position
	minDistance := math.Inf(1) // positive infinity

	for pos := range positions {
		dist := ComputeDistance(current, pos)
		if dist < minDistance {
			minDistance = dist
			closestPos = pos
		}
	}

	return closestPos
}

func isCentralPosition(pos Position, positions map[Position]bool) bool {
	// Check if the given position is in the set
	return positions[pos]
}

func ComputeDistance(pos1, pos2 Position) float64 {
	return math.Sqrt(math.Pow(pos2.X-pos1.X, 2) + math.Pow(pos2.Y-pos1.Y, 2))
}
