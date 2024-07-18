package layout

import "tkOptimizer/internal/key"

type Layout map[rune]*key.Key

func Parse(layout map[string]key.Position) Layout {
  l := make(Layout)
	for char, pos := range layout {
    rC := []rune(char)[0]
		l[rC] = key.New(pos)
	}
  return l
}
