package model

import (
	"fmt"
	"os"

	"chatbot/internal/config"
)

func (m *Model) SaveConversation(userInput, botResponse string) error {
	f, err := os.OpenFile(config.ConversationsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	info, _ := f.Stat()
	if info.Size() > 0 {
		_, err = f.WriteString("\n")
		if err != nil {
			return err
		}
	}

	_, err = f.WriteString(userInput + "|" + botResponse + "\n")
	return err
}

func (m *Model) SaveModel() error {
	fmt.Println("Saving model to", config.ModelFile)
	if err := m.Brain.Save(config.ModelFile); err != nil {
		return err
	}

	vocabFile := config.ModelFile + ".vocab"
	fmt.Println("Saving vocab to", vocabFile)
	return SaveVocab(vocabFile)
}

func (m *Model) LoadModel() error {
	if err := m.Brain.Load(config.ModelFile); err != nil {
		return err
	}

	vocabFile := config.ModelFile + ".vocab"
	if err := LoadVocab(vocabFile); err != nil {
		fmt.Println("No vocab found, will create new one")
	}
	return nil
}
