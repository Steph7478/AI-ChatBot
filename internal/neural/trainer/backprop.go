package trainer

import (
	"chatbot/internal/neural/transformer"
)

func (t *Trainer) backward(gradLogits [][]float64, input []int, cache *transformer.ForwardCache) {
	lr := t.Config.LearningRate

	var lastX [][]float64
	if len(cache.BlockOutputs) > 0 {
		lastX = cache.BlockOutputs[len(cache.BlockOutputs)-1].FFNOut
	} else {
		lastX = cache.EmbeddingOut
	}

	gradInput := t.Model.Output.Backward(gradLogits, lastX, lr)

	if len(gradInput) > 0 {
		for i := 0; i < len(input) && i < len(gradInput); i++ {
			id := input[i]
			if id >= 0 && id < len(t.Model.Embedding.Weights) {
				for j := 0; j < len(gradInput[i]) && j < len(t.Model.Embedding.Weights[id]); j++ {
					t.Model.Embedding.Weights[id][j] -= lr * gradInput[i][j]
				}
			}
		}
	}
}
