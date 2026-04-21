package main

import (
	"bufio"
	"chatbot/internal/config"
	"chatbot/internal/model"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	trainFlag := flag.Int("train", 0, "Train model for N epochs")
	flag.Parse()

	m := model.NewModel()
	m.LoadModel()

	if *trainFlag > 0 {
		fmt.Printf("🧠 Training model for %d epochs...\n", *trainFlag)
		m.Train(*trainFlag)
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	temp := config.DefaultTemp

	cmdHandler := NewCommandHandler(m, &temp)

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
