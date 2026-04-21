package trainer

import (
	"chatbot/internal/core"
	"chatbot/internal/neural/transformer"
)

type Trainer struct {
	Model  *transformer.Transformer
	Config core.TrainingConfig
}
