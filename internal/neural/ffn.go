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
	hidden := matMul(x, f.W1)
	for i := 0; i < len(hidden); i++ {
		for j := 0; j < len(hidden[i]); j++ {
			if hidden[i][j] < 0 {
				hidden[i][j] = 0
			}
		}
	}
	return matMul(hidden, f.W2)
}

func (f *FeedForwardLayer) Backward(grad [][]float64) (gradInput, gradW1, gradW2 [][]float64) {
	if len(grad) == 0 {
		return grad, nil, nil
	}

	gradW2 = matMul(transpose(grad), f.W2)

	gradHidden := matMul(grad, transpose(f.W2))
	for i := 0; i < len(gradHidden); i++ {
		for j := 0; j < len(gradHidden[i]); j++ {
			if gradHidden[i][j] < 0 {
				gradHidden[i][j] = 0
			}
		}
	}

	gradW1 = matMul(transpose(gradHidden), f.W1)
	gradInput = matMul(gradHidden, transpose(f.W1))

	return gradInput, gradW1, gradW2
}
