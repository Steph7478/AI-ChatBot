package layers

import "chatbot/internal/core"

type FeedForwardLayer struct {
	W1, W2 [][]float64
}

func NewFeedForwardLayer(inputDim, hiddenDim int) *FeedForwardLayer {
	return &FeedForwardLayer{
		W1: core.NewMatrix(inputDim, hiddenDim),
		W2: core.NewMatrix(hiddenDim, inputDim),
	}
}

func (f *FeedForwardLayer) Forward(x [][]float64) [][]float64 {
	if len(x) == 0 {
		return x
	}
	h := core.MatMul(x, f.W1)
	for i := range h {
		for j := range h[i] {
			if h[i][j] < 0 {
				h[i][j] = 0
			}
		}
	}
	return core.MatMul(h, f.W2)
}
