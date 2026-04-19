package neural

import (
	"fmt"
	"math"
)

func NewTrainer(model *Transformer, cfg TrainingConfig) *Trainer {
	return &Trainer{Model: model, Config: cfg}
}

func (t *Trainer) Train(epochs int, inputs, targets [][]int) float64 {
	if len(inputs) == 0 {
		return 0
	}
	fmt.Printf("Training: %d samples, %d epochs, LR=%.6f\n", len(inputs), epochs, t.Config.LearningRate)
	bestLoss, bestEpoch, wait := math.MaxFloat64, 0, 0
	for epoch := range epochs {
		epochLoss := 0.0
		for idx := range inputs {
			if len(inputs[idx]) == 0 || len(targets[idx]) == 0 {
				continue
			}
			logits := t.Model.Forward(inputs[idx])
			if len(logits) == 0 || len(logits[len(logits)-1]) == 0 {
				continue
			}
			last := logits[len(logits)-1]
			mx := last[0]
			for _, v := range last[1:] {
				if v > mx {
					mx = v
				}
			}
			sum := 0.0
			p := make([]float64, len(last))
			for i, v := range last {
				p[i] = math.Exp(v - mx)
				sum += p[i]
			}
			if sum == 0 {
				sum = 1e-8
			}
			for i := range p {
				p[i] /= sum
				if p[i] < 1e-10 {
					p[i] = 1e-10
				}
			}
			loss := 0.0
			for _, targ := range targets[idx] {
				if targ >= 0 && targ < len(p) {
					loss -= math.Log(p[targ])
				}
			}
			loss /= float64(len(targets[idx]))
			epochLoss += loss
			g := make([]float64, len(last))
			for i := range g {
				g[i] = p[i]
			}
			for _, targ := range targets[idx] {
				if targ >= 0 && targ < len(g) {
					g[targ] -= 1.0
				}
			}
			for i := range g {
				g[i] /= float64(len(targets[idx]))
			}
			if len(t.Model.Output.Weights) > 0 && len(t.Model.Output.Weights[0]) == len(g) {
				for i := range t.Model.Output.Weights {
					for j := range t.Model.Output.Weights[i] {
						if j < len(g) {
							t.Model.Output.Weights[i][j] -= t.Config.LearningRate * g[j]
						}
					}
				}
			}
		}
		avgLoss := epochLoss / float64(len(inputs))
		fmt.Printf("Epoch %d/%d, Loss: %.6f\n", epoch+1, epochs, avgLoss)
		if avgLoss < bestLoss {
			bestLoss, bestEpoch, wait = avgLoss, epoch+1, 0
		} else {
			wait++
		}
		if wait >= 5 && epoch+1 > 1 {
			fmt.Printf("\n🛑 EARLY STOPPING! Best loss was %.6f at epoch %d\n", bestLoss, bestEpoch)
			break
		}
	}
	return bestLoss
}
