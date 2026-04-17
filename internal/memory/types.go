package memory

import "chatbot/internal/dataset"

type Candidate struct {
	Index int
	Score float64
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

type SearchConfig struct {
	SimilarityThreshold float64
	ExactMatchBoost     float64
	MinWordMatchScore   float64
	TopKCandidates      int
	TopKTempHigh        int
	TopKTempLow         int
	IDFSmoothing        float64
}
