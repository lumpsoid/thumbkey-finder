package weights

import "errors"

type Weights [][]float64

func New(height int, width int, filler float64) Weights {
	weightsNew := make([][]float64, height)

	for i := 0; i < len(weightsNew); i++ {
		weightsNew[i] = make([]float64, width)
		for j := 0; j < len(weightsNew[0]); j++ {
			weightsNew[i][j] = filler
		}
	}
	return weightsNew
}

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
