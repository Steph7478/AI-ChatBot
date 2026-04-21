package model

import (
	"chatbot/internal/config"
	"chatbot/internal/neural/transformer"
	"fmt"
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

	m.Matcher = &SimpleTextMatcher{
		Conversations: m.Conversations,
	}
	return m
}
