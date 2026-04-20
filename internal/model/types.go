package model

import (
	"chatbot/internal/neural/transformer"
)

type ResponseResult struct {
	Text       string
	Type       ResponseType
	Confidence float64
}

type ResponseType int

const (
	ResponseGenerated ResponseType = iota
	ResponseFallback
)

type SimpleTextMatcher struct {
	Conversations map[string]string
}

type Model struct {
	Brain         *transformer.Transformer
	Matcher       *SimpleTextMatcher
	Conversations map[string]string
	Synonyms      map[string]string
}
