package keyboard

import "tkOptimizer/internal/key"

type KeyboardOption func(*Keyboard)

func applyOptions(keyboard *Keyboard, options ...KeyboardOption) {
	for _, option := range options {
		option(keyboard)
	}
}

func (k *Keyboard) Update(options ...KeyboardOption) {
	for _, option := range options {
		option(k)
	}
}

func SetWeights(weights Weights) KeyboardOption {
	return func(keyboard *Keyboard) {
		keyboard.Weights = weights
	}
}

func SetLayout(layout map[rune]key.Position) KeyboardOption {
	newLayout := make(map[rune]*key.Key)
	for char, pos := range layout {
		newLayout[char] = key.New(pos)
	}
	return func(keyboard *Keyboard) {
		keyboard.Layout = newLayout
	}
}
