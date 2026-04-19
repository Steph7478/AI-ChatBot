package neural

func NewFeedForwardLayer(inputDim, hiddenDim int) *FeedForwardLayer {
	return &FeedForwardLayer{
		W1: newMatrix(inputDim, hiddenDim),
		W2: newMatrix(hiddenDim, inputDim),
	}
}

func (f *FeedForwardLayer) Forward(x [][]float64) [][]float64 {
	if len(x) == 0 {
		return x
	}
	h := matMul(x, f.W1)
	for i := range h {
		for j := range h[i] {
			if h[i][j] < 0 {
				h[i][j] = 0
			}
		}
	}
	return matMul(h, f.W2)
}
