package app

import (
	"fmt"
	"strings"

	"chatbot/internal/dataset"
	"chatbot/internal/model"
)

func ShowStats(vocab *dataset.Vocabulary, m *model.Model) {
	fmt.Printf("\n📊 Ready | %d words | %d convs\n", vocab.Size, len(m.Dataset.Conversations))
	fmt.Println("💬 Commands: /quit, /stats, /temp N, /save, /reload")
	fmt.Println(strings.Repeat("=", 50))
}
