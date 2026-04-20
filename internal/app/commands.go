package app

import (
	"fmt"
	"os"
	"strings"

	"chatbot/internal/config"
	"chatbot/internal/model"
)

type CommandHandler struct {
	model *model.Model
	temp  *float64
}

func NewCommandHandler(m *model.Model, temp *float64) *CommandHandler {
	return &CommandHandler{model: m, temp: temp}
}

func (h *CommandHandler) Handle(input string) bool {
	switch {
	case input == "/quit":
		fmt.Println("👋 Bye!")
		os.Exit(0)

	case input == "/save":
		h.model.SaveModel()
		fmt.Println("✅ Model saved")

	case input == "/stats":
		fmt.Printf("📊 Model stats:\n")
		fmt.Printf("   Vocab size: %d\n", config.VocabSize)
		fmt.Printf("   Embedding dim: %d\n", config.EmbeddingDim)
		fmt.Printf("   Hidden dim: %d\n", config.HiddenDim)
		fmt.Printf("   Num layers: %d\n", config.NumLayers)
		fmt.Printf("   Num heads: %d\n", config.NumHeads)
		fmt.Printf("   Conversations: %d\n", len(h.model.Conversations))

	case strings.HasPrefix(input, "/temp"):
		var t float64
		fmt.Sscanf(input, "/temp %f", &t)
		*h.temp = clamp(t, config.MinTemp, config.MaxTemp)
		fmt.Printf("🌡️ Temperature: %.1f\n", *h.temp)

	case input == "/train":
		fmt.Println("🧠 Starting training with 10 epochs...")
		h.model.Train(10)

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
