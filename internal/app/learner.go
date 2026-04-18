package app

import (
	"bufio"
	"fmt"
	"strings"

	"chatbot/internal/dataset"
	"chatbot/internal/model"
)

type Learner struct {
	vocab *dataset.Vocabulary
	model *model.Model
}

func NewLearner(vocab *dataset.Vocabulary, m *model.Model) *Learner {
	return &Learner{
		vocab: vocab,
		model: m,
	}
}

func (l *Learner) LearnFromUser(scanner *bufio.Scanner, input string) {
	fmt.Print("📝 Teach me: ")
	scanner.Scan()
	correct := strings.TrimSpace(scanner.Text())
	if correct == "" {
		return
	}

	l.model.LearnFromUser(input, correct)

	convMem := model.NewConversationMemory(l.vocab)
	convMem.Load("data/checkpoint.gob")

	userTokens := l.vocab.AddAndTokenize(input)
	botTokens := l.vocab.AddAndTokenize(correct)
	convMem.Learn(userTokens, botTokens)
	convMem.CalculateIDF()
	convMem.Save("data/checkpoint.gob")

	fmt.Println("✅ Learned! Thank you ☕")
}
