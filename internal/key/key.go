package key

type GestureType bool

const (
	Press GestureType = true
	Swipe GestureType = false
)

type Key struct {
	Type     GestureType
	Central  Position
	Position Position
}

func New(pos Position) *Key {
	var keyType = GestureType(isCentralPosition(pos, mapOfCentralPositions))
	if keyType == Press {
		return &Key{keyType, pos, pos}
	}
	centralPos := closestPosition(pos, mapOfCentralPositions)
	return &Key{keyType, centralPos, pos}
}
