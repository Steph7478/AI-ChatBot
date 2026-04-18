package model

import (
	"fmt"

	"chatbot/internal/config"
	"chatbot/internal/neural"
)

func NewModel() *Model {
	m := &Model{
		Brain: neural.NewTransformer(
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

	if err := m.LoadModel(); err != nil {
		fmt.Println("No saved model found, starting fresh")
	}

	m.Matcher = NewSimpleTextMatcher(m.Conversations)

	return m
}
