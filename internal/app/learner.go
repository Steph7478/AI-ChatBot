package app

import (
	"bufio"
	"fmt"
	"strings"

	"chatbot/internal/model"
)

type Learner struct {
	model *model.Model
}

func NewLearner(m *model.Model) *Learner {
	return &Learner{model: m}
}

func (l *Learner) LearnFromUser(scanner *bufio.Scanner, input string) {
	if response, exists := l.model.Conversations[input]; exists {
		fmt.Printf("Bot: %s\n", response)
		return
	}

	normalized := l.model.Normalize(input)
	if response, exists := l.model.Conversations[normalized]; exists {
		fmt.Printf("Bot: %s\n", response)
		return
	}

	if match, score := l.model.Matcher.FindBestMatch(input); match != "" && score > 0.3 {
		fmt.Printf("Bot: %s\n", match)
		return
	}

	fmt.Print("📝 Teach me: ")
	scanner.Scan()
	teaching := strings.TrimSpace(scanner.Text())

	if teaching == "" {
		fmt.Println("❌ No teaching provided")
		return
	}

	l.model.Learn(input, teaching)
	fmt.Println("✅ Learned! 🧠")
}
