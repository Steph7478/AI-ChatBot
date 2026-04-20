package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"chatbot/internal/app"
	"chatbot/internal/config"
	"chatbot/internal/model"
)

func main() {
	trainFlag := flag.Int("train", 0, "Train model for N epochs")
	flag.Parse()

	m := model.NewModel()
	m.LoadAll()

	if *trainFlag > 0 {
		fmt.Printf("🧠 Training model for %d epochs...\n", *trainFlag)
		m.Train(*trainFlag)
		fmt.Println("✅ Training complete! Model saved.")
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	temp := config.DefaultTemp

	cmdHandler := app.NewCommandHandler(m, &temp)

	fmt.Println("🧠 Neural Chatbot")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("\n📊 Ready | Temp: %.1f\n", temp)
	fmt.Println("💬 Commands: /quit, /stats, /temp N, /save, /train")
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

		m.GenerateResponse(input, temp, scanner)
	}
}
