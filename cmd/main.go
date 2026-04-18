package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"chatbot/internal/app"
	"chatbot/internal/config"
	"chatbot/internal/model"
)

func main() {
	fmt.Println("🧠 Neural Chatbot")
	fmt.Println(strings.Repeat("=", 50))

	m := model.NewModel()
	m.LoadAll()

	scanner := bufio.NewScanner(os.Stdin)
	temp := config.DefaultTemp

	cmdHandler := app.NewCommandHandler(m, &temp)
	learner := app.NewLearner(m)

	fmt.Printf("\n📊 Ready | Temp: %.1f\n", temp)
	fmt.Println("💬 Commands: /quit, /stats, /temp N, /save")
	fmt.Println(strings.Repeat("=", 50))

	for {
		fmt.Print("\n\033[34mYou: ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		fmt.Print("\033[0m")

		if input == "" {
			continue
		}

		if cmdHandler.Handle(input) {
			continue
		}

		result := m.GenerateResponse(input, temp)
		fmt.Printf("\033[36mBot: %s\033[0m\n", result.Text)

		if result.Type == model.ResponseFallback {
			learner.LearnFromUser(scanner, input)
		}
	}
}
