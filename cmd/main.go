package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"chatbot/internal/dataset"
	"chatbot/internal/memory"
)

func main() {
	fmt.Println("🤖 Chatbot - Learning by TF-IDF Similarity")
	fmt.Println("===========================================")
	fmt.Println()

	maxTokens := 30
	temperature := 0.7
	checkpointFile := "data/memory_checkpoint.gob"

	vocab := dataset.NewVocabulary()
	data := dataset.NewDataset(vocab)

	var mem *memory.ConversationMemory

	fmt.Println("📀 Checking for saved memory...")
	mem = memory.NewConversationMemory(vocab)

	if err := mem.Load(checkpointFile); err == nil {
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

		if _, err := os.Stat("data/conversations.txt"); os.IsNotExist(err) {
			fmt.Println("❌ data/conversations.txt not found!")
			return
		}

		if err := data.LoadFromFile("data/conversations.txt"); err != nil {
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
		if err := mem.Save(checkpointFile); err != nil {
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
			mem.Save(checkpointFile)
			fmt.Println("👋 Goodbye!")
			return

		case input == "/save":
			if err := mem.Save(checkpointFile); err != nil {
				fmt.Printf("❌ Error saving: %v\n", err)
			} else {
				fmt.Println("✅ Checkpoint saved!")
			}
			continue

		case input == "/reload":
			fmt.Println("📀 Reloading from conversations.txt...")
			newVocab := dataset.NewVocabulary()
			newData := dataset.NewDataset(newVocab)

			if err := newData.LoadFromFile("data/conversations.txt"); err != nil {
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

			mem.Save(checkpointFile)
			fmt.Printf("✅ Reloaded %d conversations!\n", len(mem.Questions))
			continue

		case input == "/stats":
			fmt.Printf("\n📊 Stats:\n")
			fmt.Printf("   Conversations: %d\n", len(mem.Questions))
			fmt.Printf("   Vocabulary: %d words\n", vocab.Size)
			fmt.Printf("   Temperature: %.1f\n", temperature)
			fmt.Println()
			continue

		case strings.HasPrefix(input, "/temp "):
			fmt.Sscanf(input, "/temp %f", &temperature)
			if temperature < 0.1 {
				temperature = 0.1
			}
			if temperature > 1.5 {
				temperature = 1.5
			}
			fmt.Printf("✅ Temperature set to: %.1f\n", temperature)
			continue
		}

		questionTokens := data.Tokenize(input)
		responseTokens := mem.Generate(questionTokens, maxTokens, temperature)

		if len(responseTokens) == 0 {
			if containsAny(input, []string{"hello", "hi", "hey"}) {
				fmt.Println("Bot: hello how can I help you today")
			} else {
				fmt.Println("Bot: I don't know how to answer that yet")
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
