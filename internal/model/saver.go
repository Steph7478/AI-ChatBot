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
