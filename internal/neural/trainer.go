package neural

import (
	"fmt"
	"math"
)

func NewTrainer(model *Transformer, cfg TrainingConfig) *Trainer {
	return &Trainer{
		Model:  model,
		Config: cfg,
	}
}

func (t *Trainer) Train(epochs int, inputs, targets [][]int) float64 {
	if len(inputs) == 0 {
		return 0
	}

	fmt.Printf("Training: %d samples, %d epochs, LR=%.6f\n", len(inputs), epochs, t.Config.LearningRate)
	fmt.Printf("Initial Output weights[0][0:5]: %v\n", t.Model.Output.Weights[0][:5])

	totalLoss := 0.0
	bestLoss := math.MaxFloat64
	bestEpoch := 0
	noImproveCount := 0
	patience := 5

	for epoch := 0; epoch < epochs; epoch++ {
		epochLoss := 0.0
		gradNormSum := 0.0
		weightNormSum := 0.0

		for idx := 0; idx < len(inputs); idx++ {
			if len(inputs[idx]) == 0 || len(targets[idx]) == 0 {
				continue
			}

			weightsBefore := make([][]float64, len(t.Model.Output.Weights))
			for i := range t.Model.Output.Weights {
				weightsBefore[i] = make([]float64, len(t.Model.Output.Weights[i]))
				copy(weightsBefore[i], t.Model.Output.Weights[i])
			}

			logits := t.Model.Forward(inputs[idx])
			if len(logits) == 0 || len(logits[len(logits)-1]) == 0 {
				continue
			}

			lastLogits := logits[len(logits)-1]

			maxLogit := lastLogits[0]
			for _, v := range lastLogits {
				if v > maxLogit {
					maxLogit = v
				}
			}

			expSum := 0.0
			probs := make([]float64, len(lastLogits))
			for i, v := range lastLogits {
				probs[i] = math.Exp(v - maxLogit)
				expSum += probs[i]
			}
			if expSum == 0 {
				expSum = 1e-8
			}
			for i := range probs {
				probs[i] /= expSum
				if probs[i] < 1e-10 {
					probs[i] = 1e-10
				}
			}

			loss := 0.0
			for _, target := range targets[idx] {
				if target >= 0 && target < len(probs) {
					loss -= math.Log(probs[target])
				}
			}
			loss /= float64(len(targets[idx]))
			epochLoss += loss

			grad := make([]float64, len(lastLogits))
			for i := range grad {
				grad[i] = probs[i]
			}
			for _, target := range targets[idx] {
				if target >= 0 && target < len(grad) {
					grad[target] -= 1.0
				}
			}
			for i := range grad {
				grad[i] /= float64(len(targets[idx]))
			}

			gradNorm := 0.0
			for _, g := range grad {
				gradNorm += g * g
			}
			gradNorm = math.Sqrt(gradNorm)
			gradNormSum += gradNorm

			if len(t.Model.Output.Weights) > 0 && len(t.Model.Output.Weights[0]) == len(grad) {
				for i := 0; i < len(t.Model.Output.Weights); i++ {
					for j := 0; j < len(t.Model.Output.Weights[i]) && j < len(grad); j++ {
						update := t.Config.LearningRate * grad[j]
						t.Model.Output.Weights[i][j] -= update
					}
				}
			}

			weightChange := 0.0
			for i := range t.Model.Output.Weights {
				for j := range t.Model.Output.Weights[i] {
					diff := t.Model.Output.Weights[i][j] - weightsBefore[i][j]
					weightChange += diff * diff
				}
			}
			weightNormSum += math.Sqrt(weightChange)
		}

		avgLoss := epochLoss / float64(len(inputs))
		avgGradNorm := gradNormSum / float64(len(inputs))
		avgWeightChange := weightNormSum / float64(len(inputs))
		totalLoss = avgLoss

		if epoch == 0 {
			fmt.Printf("After epoch 1, Output weights[0][0:5]: %v\n", t.Model.Output.Weights[0][:5])
		}

		fmt.Printf("Epoch %d/%d, Loss: %.6f, GradNorm: %.6f, WeightChange: %.6f\n",
			epoch+1, epochs, avgLoss, avgGradNorm, avgWeightChange)

		if epoch > 0 && avgLoss > totalLoss && epoch == 1 {
			fmt.Println("⚠️  WARNING: Loss is increasing! Gradient direction might be wrong.")
		}

		if avgLoss < bestLoss {
			bestLoss = avgLoss
			bestEpoch = epoch + 1
			noImproveCount = 0
			fmt.Printf("  🎯 New best loss! (%.6f)\n", bestLoss)
		} else {
			noImproveCount++
			fmt.Printf("  ⏳ No improvement for %d/%d epochs\n", noImproveCount, patience)
		}

		if noImproveCount >= patience && epoch+1 > 1 {
			fmt.Printf("\n🛑 EARLY STOPPING! Best loss was %.6f at epoch %d\n", bestLoss, bestEpoch)
			break
		}
	}

	return totalLoss
}
