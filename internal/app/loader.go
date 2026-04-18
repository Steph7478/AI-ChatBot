package app

import (
	"fmt"
	"os"

	"chatbot/internal/config"
	"chatbot/internal/dataset"
	"chatbot/internal/model"
)

type Loader struct {
	vocab *dataset.Vocabulary
	model *model.Model
}

func NewLoader(vocab *dataset.Vocabulary, m *model.Model) *Loader {
	return &Loader{
		vocab: vocab,
		model: m,
	}
}

func (l *Loader) LoadOrCreate() {
	if _, err := os.Stat(config.CheckpointFile); err == nil {
		convMem := model.NewConversationMemory(l.vocab)
		if convMem.Load(config.CheckpointFile) == nil {
			fmt.Printf("📀 Loaded %d conversations from checkpoint\n", len(convMem.Questions))
			l.rebuildFromCheckpoint(convMem)
			return
		}
	}

	l.loadFromFile()
}

func (l *Loader) rebuildFromCheckpoint(convMem *model.ConversationMemory) {
	fmt.Println("🔄 Rebuilding model from checkpoint...")
	for i := 0; i < len(convMem.Questions); i++ {
		userText := convMem.QuestionText[i]
		botText := convMem.AnswerText[i]
		l.model.LearnFromConversation(userText, botText)
	}
}

func (l *Loader) loadFromFile() {
	fmt.Println("📖 Loading from conversations.txt...")
	data := dataset.NewDataset(l.vocab)

	if err := data.LoadFromFile(config.ConversationsFile); err != nil {
		fmt.Printf("⚠️ Warning: %v\n", err)
		return
	}

	fmt.Printf("✅ Loaded %d conversations from file\n", len(data.Conversations))

	convMem := model.NewConversationMemory(l.vocab)

	for _, conv := range data.Conversations {
		userText := data.Detokenize(conv[0])
		botText := data.Detokenize(conv[1])

		convMem.Learn(conv[0], conv[1])
		l.model.LearnFromConversation(userText, botText)
	}

	fmt.Println("📊 Calculating TF-IDF vectors...")
	convMem.CalculateIDF()

	fmt.Println("💾 Saving memory to checkpoint...")
	convMem.Save(config.CheckpointFile)
}
