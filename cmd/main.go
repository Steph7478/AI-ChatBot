package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"chatbot/internal/app"
	"chatbot/internal/config"
	"chatbot/internal/dataset"
	"chatbot/internal/model"
)

func main() {
	fmt.Println("☕ Java Chatbot")
	fmt.Println(strings.Repeat("=", 50))

	vocab := dataset.NewVocabulary()
	m := model.NewModel(vocab)

	loader := app.NewLoader(vocab, m)
	loader.LoadOrCreate()

	app.ShowStats(vocab, m)

	scanner := bufio.NewScanner(os.Stdin)
	temp := config.DefaultTemp
	commandHandler := app.NewCommandHandler(m, &temp)
	learner := app.NewLearner(vocab, m)

	for {
		fmt.Print("\nYou: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		input = strings.TrimPrefix(input, "You:")
		input = strings.TrimSpace(input)
		input = strings.TrimPrefix(input, "you:")
		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		if commandHandler.Handle(input) {
			continue
		}

		result := m.FindResponse(input, temp)
		fmt.Printf("Bot: %s\n", result.Text)

		if result.Type == model.ResponseFallback {
			learner.LearnFromUser(scanner, input)
		}
	}
}
