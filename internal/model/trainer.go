package model

import (
	"fmt"

	"chatbot/internal/config"
	"chatbot/internal/neural"
)

func (m *Model) Train(epochs int) {
	if len(m.Conversations) == 0 {
		fmt.Println("No conversations to train on!")
		return
	}

	fmt.Printf("Training neural network on %d conversations for %d epochs...\n", len(m.Conversations), epochs)

	inputs := make([][]int, 0)
	targets := make([][]int, 0)

	for question, answer := range m.Conversations {
		inputTokens := defaultTokenizer(question)
		outputTokens := defaultTokenizer(answer)

		if len(inputTokens) > config.TrainMaxSeqLen {
			inputTokens = inputTokens[:config.TrainMaxSeqLen]
		}
		if len(outputTokens) > config.TrainMaxSeqLen {
			outputTokens = outputTokens[:config.TrainMaxSeqLen]
		}

		if len(inputTokens) > 0 && len(outputTokens) > 0 {
			inputs = append(inputs, inputTokens)
			targets = append(targets, outputTokens)
		}
	}

	if len(inputs) == 0 {
		fmt.Println("No valid training data!")
		return
	}

	trainer := neural.NewTrainer(m.Brain, neural.TrainingConfig{
		LearningRate: config.LearningRate,
		BatchSize:    config.BatchSize,
		Epochs:       epochs,
		Patience:     3,
	})

	loss := trainer.Train(epochs, inputs, targets)
	fmt.Printf("Training completed! Final loss: %.4f\n", loss)

	m.SaveModel()
}
