package model

import (
	"bufio"
	"os"
	"strings"

	"chatbot/internal/config"
)

func (m *Model) LoadAll() {
	m.LoadPrompts()
	m.LoadExamples()
	m.LoadConversations()
}

func (m *Model) LoadPrompts() {
	file, err := os.Open(config.PromptsFile)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, "|") {
			parts := strings.Split(line, "|")
			mainPhrase := strings.ToLower(strings.TrimSpace(parts[0]))
			for i := 1; i < len(parts); i++ {
				synonym := strings.ToLower(strings.TrimSpace(parts[i]))
				m.Synonyms[synonym] = mainPhrase
			}
		}
	}
}

func (m *Model) LoadExamples() {
	file, err := os.Open(config.ExamplesFile)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, "->") {
			parts := strings.Split(line, "->")
			if len(parts) == 2 {
				userInput := strings.ToLower(strings.TrimSpace(parts[0]))
				botResponse := strings.TrimSpace(parts[1])
				m.Conversations[userInput] = botResponse
			}
		}
	}
}

func (m *Model) LoadConversations() {
	file, err := os.Open(config.ConversationsFile)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lastUser string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if after, ok := strings.CutPrefix(line, "User:"); ok {
			lastUser = strings.ToLower(strings.TrimSpace(after))
		} else if strings.HasPrefix(line, "Bot:") && lastUser != "" {
			botResponse := strings.TrimSpace(strings.TrimPrefix(line, "Bot:"))
			m.Conversations[lastUser] = botResponse
			lastUser = ""
		}
	}
}
