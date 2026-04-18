package model

import (
	"chatbot/internal/dataset"
)

type Model struct {
	Dataset *dataset.Dataset
	Vocab   *dataset.Vocabulary
	LCS     *LCSMatcher
}

type ConversationMemory struct {
	Questions      [][]int
	Answers        [][]int
	QuestionText   []string
	AnswerText     []string
	Vocab          *dataset.Vocabulary
	QuestionTF     []map[int]float64
	IDF            map[int]float64
	VocabularySize int
}

type ResponseType int

type LCSMatcher struct {
	Documents [][]string
	Answers   [][]int
}

const (
	ResponseDataset ResponseType = iota
	ResponseFallback
)

type ResponseResult struct {
	Text       string
	Type       ResponseType
	Confidence float64
}
