package model

import (
	"encoding/gob"
	"fmt"
	"os"

	"chatbot/internal/config"
)

func (m *Model) SaveModel() error {
	fmt.Println("\n💾 ===== SAVING MODEL AND VOCAB =====")

	fmt.Println("1. Saving model to", config.ModelFile)
	if err := m.Brain.Save(config.ModelFile); err != nil {
		fmt.Printf("   ❌ Error saving model: %v\n", err)
		return err
	}
	fmt.Println("   ✅ Model saved successfully!")

	vocabFile := config.ModelFile + ".vocab"
	fmt.Println("2. Saving vocab to", vocabFile)

	if err := saveVocab(vocabFile); err != nil {
		fmt.Printf("   ❌ Error saving vocab: %v\n", err)
		return err
	}

	fmt.Println("   ✅ Vocab saved successfully!")
	fmt.Println("💾 ===== SAVE COMPLETE =====")
	return nil
}

func saveVocab(path string) error {
	mu.RLock()
	defer mu.RUnlock()

	fmt.Printf("   📊 wordToID has %d entries\n", len(wordToID))
	fmt.Printf("   📊 idToWord has %d entries\n", len(idToWord))
	fmt.Printf("   📊 nextID is %d\n", nextID)

	if len(wordToID) == 0 {
		fmt.Println("   ⚠️ WARNING: wordToID is empty! Nothing to save!")
		return nil
	}

	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("   ❌ Cannot create vocab file: %v\n", err)
		return err
	}
	defer file.Close()

	data := struct {
		WordToID map[string]int
		IDToWord map[int]string
		NextID   int
	}{
		WordToID: wordToID,
		IDToWord: idToWord,
		NextID:   nextID,
	}

	if err := gob.NewEncoder(file).Encode(data); err != nil {
		fmt.Printf("   ❌ Cannot encode vocab: %v\n", err)
		return err
	}

	return nil
}

func (m *Model) SaveConversation(userInput, botResponse string) error {
	f, err := os.OpenFile(config.ConversationsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	info, _ := f.Stat()
	if info.Size() > 0 {
		f.WriteString("\n")
	}

	_, err = f.WriteString(userInput + "|" + botResponse + "\n")
	return err
}

func (m *Model) SaveTrainingData(userInput, botResponse string) error {
	f, err := os.OpenFile(config.TrainingDataFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	info, _ := f.Stat()
	if info.Size() > 0 {
		f.WriteString("\n")
	}

	_, err = f.WriteString(userInput + "|" + botResponse + "\n")
	return err
}
