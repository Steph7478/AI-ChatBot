package model

import (
	"strings"

	"chatbot/internal/config"
	"chatbot/internal/dataset"
)

func NewModel(vocab *dataset.Vocabulary) *Model {
	return &Model{
		Dataset: dataset.NewDataset(vocab),
		Vocab:   vocab,
		LCS:     NewLCSMatcher(),
	}
}

func (m *Model) LearnFromConversation(userMsg, botMsg string) {
	m.Dataset.AddConversation(userMsg, botMsg)
	tokens := m.Vocab.Tokenize(botMsg)
	m.LCS.AddDocument(userMsg, tokens)
}

func (m *Model) LearnFromUser(userMsg, botMsg string) {
	m.Dataset.AddConversation(userMsg, botMsg)
	tokens := m.Vocab.Tokenize(botMsg)
	m.LCS.AddDocument(userMsg, tokens)
}

func (m *Model) FindResponse(question string, temp float64) ResponseResult {
	if resp := m.findExactMatch(question); resp != "" {
		return ResponseResult{resp, ResponseDataset, 0.95}
	}

	conf := m.LCS.GetConfidence(question)
	if conf > config.LCSConfidenceThreshold {
		if resp := m.LCS.Predict(question); resp != nil {
			return ResponseResult{m.Vocab.Detokenize(resp), ResponseDataset, conf}
		}
	}

	return ResponseResult{"i am still learning please teach me more conversations ☕", ResponseFallback, 0.3}
}

func (m *Model) findExactMatch(q string) string {
	qLower := strings.ToLower(q)

	for _, conv := range m.Dataset.Conversations {
		questionText := strings.ToLower(m.Dataset.Detokenize(conv[0]))
		if questionText == qLower {
			return m.Dataset.Detokenize(conv[1])
		}
	}
	return ""
}
