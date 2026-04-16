package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"chatbot/internal/config"
	"chatbot/internal/dataset"
	"chatbot/internal/memory"
)

func main() {
	fmt.Println("🤖 Chatbot - Learning by TF-IDF Similarity")
	fmt.Println("===========================================")
	fmt.Println()

	vocab := dataset.NewVocabulary()
	data := dataset.NewDataset(vocab)

	var mem *memory.ConversationMemory

	fmt.Println("📀 Checking for saved memory...")
	mem = memory.NewConversationMemory(vocab)

	if err := mem.Load(config.CheckpointFile); err == nil {
		fmt.Println("✅ Loaded memory from checkpoint!")

		for word, id := range mem.Vocab.WordToID {
			vocab.WordToID[word] = id
			vocab.IDToWord[id] = word
			if id >= vocab.Size {
				vocab.Size = id + 1
			}
		}
		data.Vocab = vocab

		fmt.Printf("✅ %d conversations loaded from checkpoint\n", len(mem.Questions))
		fmt.Printf("✅ %d words in vocabulary\n", vocab.Size)
	} else {
		fmt.Println("📝 No checkpoint found. Loading from conversations.txt...")

		if _, err := os.Stat(config.ConversationsFile); os.IsNotExist(err) {
			fmt.Printf("❌ %s not found!\n", config.ConversationsFile)
			return
		}

		if err := data.LoadFromFile(config.ConversationsFile); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		fmt.Printf("✅ %d conversations loaded from file\n", len(data.Conversations))
		fmt.Printf("✅ %d words in vocabulary\n", data.Vocab.Size)

		mem = memory.NewConversationMemory(vocab)

		for _, conv := range data.Conversations {
			mem.Learn(conv[0], conv[1])
		}

		fmt.Println("📊 Calculating TF-IDF vectors...")
		mem.CalculateIDF()

		fmt.Println("💾 Saving memory to checkpoint...")
		if err := mem.Save(config.CheckpointFile); err != nil {
			fmt.Printf("⚠️ Warning: Could not save checkpoint: %v\n", err)
		} else {
			fmt.Println("✅ Checkpoint saved!")
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("💬 Ready! I learn by TF-IDF similarity")
	fmt.Println("Commands: /quit, /stats, /temp [0.1-1.5], /save, /reload")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	temperature := config.DefaultTemperature

	for {
		fmt.Print("You: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if input == "" {
			continue
		}

		switch {
		case input == "/quit":
			fmt.Println("💾 Saving checkpoint...")
			mem.Save(config.CheckpointFile)
			fmt.Println("👋 Goodbye!")
			return

		case input == "/save":
			if err := mem.Save(config.CheckpointFile); err != nil {
				fmt.Printf("❌ Error saving: %v\n", err)
			} else {
				fmt.Println("✅ Checkpoint saved!")
			}
			continue

		case input == "/reload":
			fmt.Println("📀 Reloading from conversations.txt...")
			newVocab := dataset.NewVocabulary()
			newData := dataset.NewDataset(newVocab)

			if err := newData.LoadFromFile(config.ConversationsFile); err != nil {
				fmt.Printf("❌ Error: %v\n", err)
				continue
			}

			newMem := memory.NewConversationMemory(newVocab)
			for _, conv := range newData.Conversations {
				newMem.Learn(conv[0], conv[1])
			}
			newMem.CalculateIDF()

			mem = newMem
			data = newData
			vocab = newVocab

			mem.Save(config.CheckpointFile)
			fmt.Printf("✅ Reloaded %d conversations!\n", len(mem.Questions))
			continue

		case input == "/stats":
			fmt.Printf("\n📊 Stats:\n")
			fmt.Printf("   Conversations: %d\n", len(mem.Questions))
			fmt.Printf("   Vocabulary: %d words\n", vocab.Size)
			fmt.Printf("   Temperature: %.1f\n", temperature)
			fmt.Printf("   Similarity Threshold: %.2f\n", config.SimilarityThreshold)
			fmt.Println()
			continue

		case strings.HasPrefix(input, "/temp "):
			fmt.Sscanf(input, "/temp %f", &temperature)
			if temperature < config.MinTemperature {
				temperature = config.MinTemperature
			}
			if temperature > config.MaxTemperature {
				temperature = config.MaxTemperature
			}
			fmt.Printf("✅ Temperature set to: %.1f\n", temperature)
			continue
		}

		questionTokens := data.Tokenize(input)
		responseTokens := mem.Generate(questionTokens, config.MaxTokens, temperature)

		if len(responseTokens) == 0 {
			if containsAny(input, config.GreetingWords) {
				fmt.Printf("Bot: %s\n", config.GreetingResponse)
			} else {
				fmt.Printf("Bot: %s\n", config.DefaultResponse)
			}
		} else {
			response := data.Detokenize(responseTokens)
			fmt.Printf("Bot: %s\n", response)
		}
		fmt.Println()
	}
}

func containsAny(text string, words []string) bool {
	for _, w := range words {
		if strings.Contains(text, w) {
			return true
		}
	}
	return false
}
