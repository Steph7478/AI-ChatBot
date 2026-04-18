package neural

import (
	"fmt"
	"math"

	"chatbot/internal/config"
)

func (a *Adam) Update(name string, weights, grad [][]float64) [][]float64 {
	a.Step++

	if _, ok := a.M[name]; !ok {
		a.M[name] = newMatrix(len(weights), len(weights[0]))
		a.V[name] = newMatrix(len(weights), len(weights[0]))
	}

	updated := make([][]float64, len(weights))
	for i := range weights {
		updated[i] = make([]float64, len(weights[i]))
		for j := range weights[i] {
			a.M[name][i][j] = a.Beta1*a.M[name][i][j] + (1-a.Beta1)*grad[i][j]
			a.V[name][i][j] = a.Beta2*a.V[name][i][j] + (1-a.Beta2)*grad[i][j]*grad[i][j]

			mHat := a.M[name][i][j] / (1 - math.Pow(a.Beta1, float64(a.Step)))
			vHat := a.V[name][i][j] / (1 - math.Pow(a.Beta2, float64(a.Step)))

			update := a.LearningRate * mHat / (math.Sqrt(vHat) + a.Epsilon)
			if update > config.GradientClip {
				update = config.GradientClip
			}
			if update < -config.GradientClip {
				update = -config.GradientClip
			}
			updated[i][j] = weights[i][j] - update
		}
	}
	return updated
}

func NewTrainer(model *Transformer, cfg TrainingConfig) *Trainer {
	return &Trainer{
		Model:     model,
		Config:    cfg,
		Optimizer: NewAdam(cfg.LearningRate),
	}
}

func (t *Trainer) Train(epochs int, inputs, targets [][]int) float64 {
	if len(inputs) == 0 {
		return 0
	}

	var totalLoss float64
	lossHistory := make([]float64, epochs)

	for epoch := range epochs {
		totalLoss = 0
		validBatches := 0

		for idx := range inputs {
			logits := t.Model.Forward(inputs[idx])
			loss := crossEntropyLoss(logits, targets[idx])

			if !math.IsNaN(loss) && !math.IsInf(loss, 0) {
				totalLoss += loss
				validBatches++
				grads := t.backward(logits, targets[idx])
				t.updateWeights(grads)
			}
		}

		if validBatches > 0 {
			lossHistory[epoch] = totalLoss / float64(validBatches)
		} else {
			lossHistory[epoch] = 0
		}

		if epoch%5 == 0 || epoch == epochs-1 {
			fmt.Printf("Epoch %d/%d, Loss: %.6f\n", epoch, epochs, lossHistory[epoch])
		}
	}

	if len(lossHistory) > 0 {
		return lossHistory[len(lossHistory)-1]
	}
	return 0
}

func (t *Trainer) backward(logits [][]float64, targets []int) map[string][][]float64 {
	grads := make(map[string][][]float64)

	if len(logits) == 0 || t.Model.Output.OutputCache == nil {
		return grads
	}

	outputGrad := softmaxGradient(logits, targets)

	if len(outputGrad) == 0 {
		return grads
	}

	if len(t.Model.Output.OutputCache) > 0 && len(outputGrad) > 0 {
		outputWeightGrad := matMul(transpose(t.Model.Output.OutputCache), outputGrad)
		grads["output_weights"] = outputWeightGrad
	}

	grad := outputGrad

	for i := len(t.Model.Blocks) - 1; i >= 0; i-- {
		if i >= len(t.Model.Blocks) {
			continue
		}
		block := t.Model.Blocks[i]
		if block.FFN != nil {
			grad = block.FFN.Backward(grad, grads, i)
		}
		if block.Attention != nil {
			grad = block.Attention.Backward(grad, grads, i)
		}
	}

	if len(grad) > 0 && len(t.Model.Embedding.Weights) > 0 {
		embeddingGrad := matMul(grad, transpose(t.Model.Embedding.Weights))
		grads["embedding"] = embeddingGrad
	}

	return grads
}

func (t *Trainer) updateWeights(grads map[string][][]float64) {
	if grad, ok := grads["output_weights"]; ok {
		t.Model.Output.Weights = t.Optimizer.Update("output", t.Model.Output.Weights, grad)
	}

	for i := range t.Model.Blocks {
		if grad, ok := grads[ffnKey(i, "w1")]; ok {
			t.Model.Blocks[i].FFN.W1 = t.Optimizer.Update(ffnKey(i, "w1"), t.Model.Blocks[i].FFN.W1, grad)
		}
		if grad, ok := grads[ffnKey(i, "w2")]; ok {
			t.Model.Blocks[i].FFN.W2 = t.Optimizer.Update(ffnKey(i, "w2"), t.Model.Blocks[i].FFN.W2, grad)
		}
		if grad, ok := grads[attnKey(i, "wq")]; ok {
			t.Model.Blocks[i].Attention.WQ = t.Optimizer.Update(attnKey(i, "wq"), t.Model.Blocks[i].Attention.WQ, grad)
		}
		if grad, ok := grads[attnKey(i, "wk")]; ok {
			t.Model.Blocks[i].Attention.WK = t.Optimizer.Update(attnKey(i, "wk"), t.Model.Blocks[i].Attention.WK, grad)
		}
		if grad, ok := grads[attnKey(i, "wv")]; ok {
			t.Model.Blocks[i].Attention.WV = t.Optimizer.Update(attnKey(i, "wv"), t.Model.Blocks[i].Attention.WV, grad)
		}
		if grad, ok := grads[attnKey(i, "wo")]; ok {
			t.Model.Blocks[i].Attention.WO = t.Optimizer.Update(attnKey(i, "wo"), t.Model.Blocks[i].Attention.WO, grad)
		}
	}

	if grad, ok := grads["embedding"]; ok {
		t.Model.Embedding.Weights = t.Optimizer.Update("embedding", t.Model.Embedding.Weights, grad)
	}
}

func crossEntropyLoss(logits [][]float64, targets []int) float64 {
	if len(logits) == 0 {
		return 0
	}

	lastLogits := logits[len(logits)-1]
	probs := softmax(lastLogits)

	loss := 0.0
	count := 0
	for _, target := range targets {
		if target >= 0 && target < len(probs) {
			loss -= math.Log(probs[target] + config.Epsilon)
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return loss / float64(count)
}

func softmaxGradient(logits [][]float64, targets []int) [][]float64 {
	if len(logits) == 0 {
		return [][]float64{}
	}

	lastLogits := logits[len(logits)-1]
	probs := softmax(lastLogits)

	grad := make([][]float64, len(logits))
	for i := range grad {
		grad[i] = make([]float64, len(lastLogits))
		if i == len(logits)-1 {
			for j := range grad[i] {
				grad[i][j] = probs[j]
			}
			for _, target := range targets {
				if target >= 0 && target < len(grad[i]) {
					grad[i][target] -= 1.0
				}
			}
		}
	}
	return grad
}
