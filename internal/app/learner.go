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
	fmt.Print("📝 Teach me: ")
	scanner.Scan()
	correct := strings.TrimSpace(scanner.Text())

	if correct == "" {
		fmt.Println("❌ No teaching provided")
		return
	}

	if err := l.model.LearnAndSave(input, correct); err != nil {
		fmt.Printf("⚠️ Saved to memory but couldn't save to file: %v\n", err)
	} else {
		fmt.Printf("✅ Learned! Response for '%s' is '%s'\n", input, correct)
	}
}
