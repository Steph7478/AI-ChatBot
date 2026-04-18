package model

import "chatbot/internal/neural"

type ResponseResult struct {
	Text       string
	Type       ResponseType
	Confidence float64
}

type SimpleTextMatcher struct {
	Conversations map[string]string
}

type ResponseType int

const (
	ResponseGenerated ResponseType = iota
	ResponseFallback
)

type Model struct {
	Brain         *neural.Transformer
	Matcher       *SimpleTextMatcher
	Conversations map[string]string
	Synonyms      map[string]string
}
