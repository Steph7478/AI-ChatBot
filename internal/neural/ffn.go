package neural

type FeedForwardLayer struct {
	W1, W2 [][]float64
}

func NewFeedForwardLayer(inputDim, hiddenDim int) *FeedForwardLayer {
	return &FeedForwardLayer{
		W1: newMatrix(inputDim, hiddenDim),
		W2: newMatrix(hiddenDim, inputDim),
	}
}

func (f *FeedForwardLayer) Forward(x [][]float64) [][]float64 {
	hidden := matMul(x, f.W1)
	for i := range hidden {
		for j := range hidden[i] {
			if hidden[i][j] < 0 {
				hidden[i][j] = 0
			}
		}
	}
	return matMul(hidden, f.W2)
}
