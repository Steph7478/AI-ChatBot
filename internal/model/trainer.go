package model

import (
	"fmt"
	"os"

	"chatbot/internal/config"
	"chatbot/internal/core"
	"chatbot/internal/neural/trainer"
)

func (m *Model) Train(epochs int) {
	if len(m.TrainingData) == 0 {
		fmt.Println("No training data found! Please add data to training_data.txt")
		return
	}

	fmt.Printf("Training on %d conversations from training_data.txt for %d epochs...\n", len(m.TrainingData), epochs)

	inputs, targets := m.prepareTrainingData()
	if len(inputs) == 0 {
		fmt.Println("No valid training data!")
		return
	}

	t := trainer.New(m.Brain, core.TrainingConfig{
		LearningRate: config.LearningRate,
		BatchSize:    config.BatchSize,
		Epochs:       epochs,
		Patience:     config.Patience,
	})

	loss := t.Train(epochs, inputs, targets)
	fmt.Printf("Training completed! Final loss: %.4f\n", loss)

	m.saveTrainingResults()
}

func (m *Model) prepareTrainingData() ([][]int, [][]int) {
	inputs := make([][]int, 0)
	targets := make([][]int, 0)

	for question, answer := range m.TrainingData {
		inputTokens := DefaultTokenizer(question)
		outputTokens := DefaultTokenizer(answer)

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

	return inputs, targets
}

func (m *Model) saveTrainingResults() {
	if err := m.SaveModel(); err != nil {
		fmt.Printf("❌ Error saving: %v\n", err)
		return
	}

	fmt.Println("✅ Model saved successfully!")
	if info, err := os.Stat(config.ModelFile); err == nil {
		fmt.Printf("   File size: %d bytes\n", info.Size())
		fmt.Printf("   Modified: %s\n", info.ModTime())
	}
}
