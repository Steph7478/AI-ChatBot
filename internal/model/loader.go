package model

import (
	"bufio"
	"os"
	"strings"

	"chatbot/internal/config"
)

func (m *Model) LoadAll() {
	m.loadCSVLike(config.ConversationsFile, m.onPair)
	m.loadCSVLike(config.ExamplesFile, m.onPair)
	m.loadCSVLike(config.PromptsFile, m.onPrompt)
}

func (m *Model) loadCSVLike(path string, callback func([]string)) {
	file, err := os.Open(path)
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

		fields := strings.FieldsFunc(line, func(r rune) bool {
			return r == '|' || r == ',' || r == ';' || r == '\t' || r == ':'
		})

		if len(fields) >= 2 {
			callback(fields)
		}
	}
}

func (m *Model) onPair(fields []string) {
	input := strings.ToLower(strings.TrimSpace(fields[0]))
	response := strings.TrimSpace(fields[1])
	if input != "" && response != "" {
		m.Conversations[input] = response
	}
}

func (m *Model) onPrompt(fields []string) {
	category := strings.TrimSpace(fields[0])
	for i := 1; i < len(fields); i++ {
		word := strings.ToLower(strings.TrimSpace(fields[i]))
		if word != "" {
			m.Synonyms[word] = category
		}
	}
}
