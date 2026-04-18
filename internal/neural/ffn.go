package neural

func NewFeedForwardLayer(inputDim, hiddenDim int) *FeedForwardLayer {
	return &FeedForwardLayer{
		W1: newMatrix(inputDim, hiddenDim),
		W2: newMatrix(hiddenDim, inputDim),
	}
}

func (f *FeedForwardLayer) Forward(x [][]float64) [][]float64 {
	if len(x) == 0 {
		return [][]float64{}
	}
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

func (f *FeedForwardLayer) Backward(grad [][]float64, grads map[string][][]float64, layerIdx int) [][]float64 {
	if len(grad) == 0 || len(f.W1) == 0 || len(f.W2) == 0 {
		return grad
	}

	w2Grad := matMul(transpose(grad), f.W2)
	if len(w2Grad) > 0 && len(w2Grad[0]) > 0 {
		grads[ffnKey(layerIdx, "w2")] = w2Grad
	}

	reluGrad := make([][]float64, len(grad))
	for i := range grad {
		reluGrad[i] = make([]float64, len(grad[i]))
		for j := range grad[i] {
			if i < len(f.W1) && j < len(f.W1[i]) && f.W1[i][j] > 0 {
				reluGrad[i][j] = grad[i][j]
			}
		}
	}

	w1Grad := matMul(transpose(reluGrad), f.W1)
	if len(w1Grad) > 0 && len(w1Grad[0]) > 0 {
		grads[ffnKey(layerIdx, "w1")] = w1Grad
	}

	return matMul(reluGrad, transpose(f.W1))
}

func ffnKey(layer int, name string) string {
	return "ffn_" + name + "_" + string(rune(layer+48))
}
