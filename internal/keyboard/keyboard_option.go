package keyboard

import "tkOptimizer/internal/layout"


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

func SetLayout(layout layout.Layout) KeyboardOption {
	return func(keyboard *Keyboard) {
		keyboard.Layout = layout
	}
}
