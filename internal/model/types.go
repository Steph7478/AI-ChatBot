package model

import "chatbot/internal/neural/transformer"

type Model struct {
	Brain         *transformer.Transformer
	Matcher       *SimpleTextMatcher
	Conversations map[string]string
	TrainingData  map[string]string
	Synonyms      map[string]string
}

type SimpleTextMatcher struct {
	Conversations map[string]string
}

type ResponseResult struct {
	Text       string
	IsFallback bool
	Confidence float64
}
