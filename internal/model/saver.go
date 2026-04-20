package model

import (
	"encoding/gob"
	"fmt"
	"os"

	"chatbot/internal/config"
)

func (m *Model) SaveModel() error {
	fmt.Println("Saving model to", config.ModelFile)
	if err := m.Brain.Save(config.ModelFile); err != nil {
		return err
	}

	vocabFile := config.ModelFile + ".vocab"
	fmt.Println("Saving vocab to", vocabFile)
	return saveVocab(vocabFile)
}

func saveVocab(path string) error {
	mu.RLock()
	defer mu.RUnlock()

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	data := struct {
		WordToID map[string]int
		IDToWord map[int]string
		NextID   int
	}{
		WordToID: wordToID,
		IDToWord: idToWord,
		NextID:   nextID,
	}

	return gob.NewEncoder(file).Encode(data)
}

func (m *Model) SaveConversation(userInput, botResponse string) error {
	f, err := os.OpenFile(config.ConversationsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	info, _ := f.Stat()
	if info.Size() > 0 {
		f.WriteString("\n")
	}

	_, err = f.WriteString(userInput + "|" + botResponse + "\n")
	return err
}

func (m *Model) SaveTrainingData(userInput, botResponse string) error {
	f, err := os.OpenFile(config.TrainingDataFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	info, _ := f.Stat()
	if info.Size() > 0 {
		f.WriteString("\n")
	}

	_, err = f.WriteString(userInput + "|" + botResponse + "\n")
	return err
}
