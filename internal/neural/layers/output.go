package layers

import "chatbot/internal/core"

func NewOutputLayer(inputDim, outputDim int) *OutputLayer {
	return &OutputLayer{
		Weights: core.NewMatrix(inputDim, outputDim),
	}
}

func (o *OutputLayer) Forward(x [][]float64) [][]float64 {
	return core.MatMul(x, o.Weights)
}

func (o *OutputLayer) Backward(gradLogits [][]float64, lastX [][]float64, lr float64) [][]float64 {
	for i := range o.Weights {
		for j := range o.Weights[i] {
			grad := 0.0
			for k := 0; k < len(gradLogits) && k < len(lastX); k++ {
				if j < len(gradLogits[k]) {
					grad += gradLogits[k][j] * lastX[k][i]
				}
			}
			o.Weights[i][j] -= lr * grad
		}
	}

	gradInput := make([][]float64, len(gradLogits))
	for i := range gradLogits {
		gradInput[i] = make([]float64, len(o.Weights))
		for j := range o.Weights {
			for k := range gradLogits[i] {
				if k < len(o.Weights[j]) {
					gradInput[i][j] += gradLogits[i][k] * o.Weights[j][k]
				}
			}
		}
	}
	return gradInput
}
