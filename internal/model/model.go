package model

import (
	"fmt"

	"chatbot/internal/config"
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
		TrainingData:  make(map[string]string),
		Synonyms:      make(map[string]string),
	}

	m.LoadAll()

	fmt.Printf("Looking for model at: %s\n", config.ModelFile)
	if err := m.LoadModel(); err != nil {
		fmt.Printf("No saved model found, starting fresh: %v\n", err)
	} else {
		fmt.Println("✅ Model loaded successfully!")
	}

	m.Matcher = NewSimpleTextMatcher(m.Conversations)
	return m
}
