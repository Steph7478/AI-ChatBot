package model

import (
	"os"

	"chatbot/internal/config"
)

func (m *Model) SaveConversation(userInput, botResponse string) error {
	f, err := os.OpenFile(config.ConversationsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(userInput + "|" + botResponse + "\n")
	return err
}

func (m *Model) LearnAndSave(input, response string) error {
	m.Learn(input, response)
	return m.SaveConversation(input, response)
}
