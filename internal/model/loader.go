package model

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"os"
	"strings"

	"chatbot/internal/config"
)

func (m *Model) LoadAll() {
	m.loadFile(config.ConversationsFile, m.onPair)
	m.loadFile(config.TrainingDataFile, m.onTraining)
	m.loadFile(config.PromptsFile, m.onPrompt)
}

func (m *Model) loadFile(path string, callback func([]string)) {
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

		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			continue
		}

		input := strings.ToLower(strings.TrimSpace(parts[0]))
		response := strings.TrimSpace(parts[1])

		if input != "" && response != "" {
			callback([]string{input, response})
		}
	}
}

func loadVocab(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var data struct {
		WordToID map[string]int
		IDToWord map[int]string
		NextID   int
	}

	if err := gob.NewDecoder(file).Decode(&data); err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()

	wordToID = data.WordToID
	idToWord = data.IDToWord
	nextID = data.NextID

	return nil
}

func (m *Model) LoadModel() error {
	if err := m.Brain.Load(config.ModelFile); err != nil {
		return err
	}

	vocabFile := config.ModelFile + ".vocab"
	if err := loadVocab(vocabFile); err != nil {
		fmt.Println("No vocab found, will create new one")
	}
	return nil
}

func (m *Model) onPair(fields []string) {
	m.Conversations[fields[0]] = fields[1]
}

func (m *Model) onTraining(fields []string) {
	m.TrainingData[fields[0]] = fields[1]
}

func (m *Model) onPrompt(fields []string) {
	if len(fields) < 2 {
		return
	}
	mainPhrase := strings.ToLower(strings.TrimSpace(fields[0]))
	synonyms := strings.SplitSeq(fields[1], "|")
	for synonym := range synonyms {
		synonym = strings.ToLower(strings.TrimSpace(synonym))
		if synonym != "" {
			m.Synonyms[synonym] = mainPhrase
		}
	}
}
