package model

import (
	"fmt"

	"chatbot/internal/config"
)

func (m *Model) SaveModel() error {
	fmt.Println("Saving model to", config.ModelFile)
	return m.Brain.Save(config.ModelFile)
}

func (m *Model) LoadModel() error {
	return m.Brain.Load(config.ModelFile)
}
