package memory

import (
	"chatbot/internal/dataset"
)

func NewConversationMemory(vocab *dataset.Vocabulary) *ConversationMemory {
	return &ConversationMemory{
		Questions:      make([][]int, 0),
		Answers:        make([][]int, 0),
		QuestionText:   make([]string, 0),
		AnswerText:     make([]string, 0),
		Vocab:          vocab,
		QuestionTF:     make([]map[int]float64, 0),
		IDF:            make(map[int]float64),
		VocabularySize: vocab.Size,
	}
}

func (m *ConversationMemory) Learn(questionTokens, answerTokens []int) {
	m.Questions = append(m.Questions, questionTokens)
	m.Answers = append(m.Answers, answerTokens)

	m.QuestionText = append(m.QuestionText, m.Vocab.Detokenize(questionTokens))
	m.AnswerText = append(m.AnswerText, m.Vocab.Detokenize(answerTokens))

	tf := m.calculateTF(questionTokens)
	m.QuestionTF = append(m.QuestionTF, tf)
}

func (m *ConversationMemory) GetQuestionCount() int {
	return len(m.Questions)
}

func (m *ConversationMemory) GetAnswerCount() int {
	return len(m.Answers)
}

func (m *ConversationMemory) GetConversation(index int) ([]int, []int, string, string, bool) {
	if index < 0 || index >= len(m.Questions) {
		return nil, nil, "", "", false
	}
	return m.Questions[index], m.Answers[index], m.QuestionText[index], m.AnswerText[index], true
}

func (m *ConversationMemory) Clear() {
	m.Questions = make([][]int, 0)
	m.Answers = make([][]int, 0)
	m.QuestionText = make([]string, 0)
	m.AnswerText = make([]string, 0)
	m.QuestionTF = make([]map[int]float64, 0)
	m.IDF = make(map[int]float64)
}
