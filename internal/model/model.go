package model

import (
	"fmt"
	"os"

	"chatbot/internal/config"
	"chatbot/internal/core"
	"chatbot/internal/neural/trainer"
	"chatbot/internal/neural/transformer"
)

func NewModel() *Model {
	m := &Model{
		Brain: transformer.New(
			config.VocabSize,
			config.EmbeddingDim,
			config.HiddenDim,
			config.NumHeads,
			config.NumLayers,
			config.MaxSeqLen,
			config.Dropout,
		),
		Conversations: make(map[string]string),
		Synonyms:      make(map[string]string),
	}

	m.LoadAll()

	fmt.Printf("Looking for model at: %s\n", config.ModelFile)

	if err := m.Brain.Load(config.ModelFile); err != nil {
		fmt.Printf("No saved model found, starting fresh: %v\n", err)
	} else {
		fmt.Println("✅ Model loaded successfully! Continuing from previous training.")
	}

	m.Matcher = NewSimpleTextMatcher(m.Conversations)

	return m
}

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

	t := trainer.New(m.Brain, core.TrainingConfig{
		LearningRate: config.LearningRate,
		BatchSize:    config.BatchSize,
		Epochs:       epochs,
		Patience:     3,
	})

	loss := t.Train(epochs, inputs, targets)
	fmt.Printf("Training completed! Final loss: %.4f\n", loss)

	fmt.Println("\n📊 Weight stats BEFORE saving:")
	m.Brain.PrintWeightStats()

	fmt.Printf("💾 Saving model to %s...\n", config.ModelFile)
	if err := m.Brain.Save(config.ModelFile); err != nil {
		fmt.Printf("❌ Error saving model: %v\n", err)
	} else {
		fmt.Println("✅ Model saved successfully!")

		if info, err := os.Stat(config.ModelFile); err == nil {
			fmt.Printf("   File size: %d bytes\n", info.Size())
			fmt.Printf("   Modified: %s\n", info.ModTime())
		}
	}
}
