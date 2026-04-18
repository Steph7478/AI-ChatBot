package model

import (
	"chatbot/internal/dataset"
)

type Model struct {
	Dataset *dataset.Dataset
	Vocab   *dataset.Vocabulary
	LCS     *LCSMatcher
}

type ResponseType int

const (
	ResponseDataset ResponseType = iota
	ResponseFallback
)

type ResponseResult struct {
	Text       string
	Type       ResponseType
	Confidence float64
}
