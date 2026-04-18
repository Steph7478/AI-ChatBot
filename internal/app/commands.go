package app

import (
	"fmt"
	"os"
	"strings"

	"chatbot/internal/config"
	"chatbot/internal/dataset"
	"chatbot/internal/model"
)

type CommandHandler struct {
	model *model.Model
	temp  *float64
}

func NewCommandHandler(m *model.Model, temp *float64) *CommandHandler {
	return &CommandHandler{
		model: m,
		temp:  temp,
	}
}

func (h *CommandHandler) Handle(input string) bool {
	switch {
	case input == "/quit":
		fmt.Println("👋 Bye!")
		os.Exit(0)

	case input == "/save":
		convMem := model.NewConversationMemory(h.model.Vocab)
		convMem.Load(config.CheckpointFile)
		convMem.Save(config.CheckpointFile)
		fmt.Println("✅ Saved")

	case input == "/stats":
		convMem := model.NewConversationMemory(h.model.Vocab)
		convMem.Load(config.CheckpointFile)
		fmt.Printf("📊 Conversations: %d\n", len(convMem.Questions))

	case strings.HasPrefix(input, "/temp"):
		fmt.Sscanf(input, "/temp %f", h.temp)
		*h.temp = clamp(*h.temp, config.MinTemp, config.MaxTemp)
		fmt.Printf("🌡️ Temperature: %.1f\n", *h.temp)

	case input == "/reload":
		fmt.Println("🔄 Reloading...")
		newVocab := dataset.NewVocabulary()
		newData := dataset.NewDataset(newVocab)
		if err := newData.LoadFromFile(config.ConversationsFile); err == nil {
			newConvMem := model.NewConversationMemory(newVocab)
			for _, conv := range newData.Conversations {
				newConvMem.Learn(conv[0], conv[1])
			}
			newConvMem.CalculateIDF()
			newConvMem.Save(config.CheckpointFile)
			fmt.Println("✅ Reloaded")
		}

	default:
		return false
	}
	return true
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
