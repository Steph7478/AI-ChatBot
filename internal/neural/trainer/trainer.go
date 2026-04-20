package trainer

import (
	"fmt"
	"math"
	"math/rand"

	"chatbot/internal/core"
	"chatbot/internal/neural/transformer"
)

type Trainer struct {
	Model  *transformer.Transformer
	Config core.TrainingConfig
}

func New(model *transformer.Transformer, cfg core.TrainingConfig) *Trainer {
	return &Trainer{Model: model, Config: cfg}
}

func (t *Trainer) Train(epochs int, inputs, targets [][]int) float64 {
	if len(inputs) == 0 {
		return 0
	}

	fmt.Printf("Training: %d samples, %d epochs, LR=%.6f\n", len(inputs), epochs, t.Config.LearningRate)

	bestLoss := math.MaxFloat64
	wait := 0

	for epoch := range epochs {
		epochLoss := t.trainEpoch(inputs, targets)

		avgLoss := epochLoss / float64(len(inputs))
		fmt.Printf("Epoch %d/%d, Loss: %.6f\n", epoch+1, epochs, avgLoss)

		if avgLoss < bestLoss {
			bestLoss = avgLoss
			wait = 0
		} else {
			wait++
		}

		if wait >= t.Config.Patience && epoch+1 > 5 {
			fmt.Printf("\n🛑 Early stopping! Best loss: %.6f\n", bestLoss)
			break
		}
	}

	return bestLoss
}

func (t *Trainer) trainEpoch(inputs, targets [][]int) float64 {
	epochLoss := 0.0
	indices := rand.Perm(len(inputs))

	for _, idx := range indices {
		input := inputs[idx]
		target := targets[idx]

		if len(input) == 0 || len(target) == 0 {
			continue
		}

		seqLen := len(target)
		if len(input) > seqLen {
			input = input[:seqLen]
		}

		logits, cache := t.Model.ForwardWithCache(input)
		if len(logits) == 0 {
			continue
		}

		loss, gradLogits := core.CrossEntropyLoss(logits, target)
		epochLoss += loss

		t.backward(gradLogits, input, cache)
	}

	return epochLoss
}
