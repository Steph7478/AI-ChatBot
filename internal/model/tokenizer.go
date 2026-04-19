package model

import (
	"encoding/gob"
	"os"
	"strings"
	"sync"

	"chatbot/internal/neural"
)

var (
	wordToID = make(map[string]int)
	idToWord = make(map[int]string)
	nextID   = 2
	mu       sync.RWMutex
)

func init() {
	wordToID["<PAD>"] = 0
	idToWord[0] = "<PAD>"
	wordToID["<UNK>"] = 1
	idToWord[1] = "<UNK>"
}

func getWordID(word string) int {
	mu.Lock()
	defer mu.Unlock()
	if id, exists := wordToID[word]; exists {
		return id
	}
	id := nextID
	wordToID[word] = id
	idToWord[id] = word
	nextID++
	return id
}

func defaultTokenizer(text string) []int {
	if text == "" {
		return []int{0}
	}

	text = strings.ToLower(text)
	words := strings.Fields(text)
	tokens := make([]int, 0, len(words))

	for _, word := range words {
		word = strings.Trim(word, ".,!?;:()[]{}\"")
		tokens = append(tokens, getWordID(word))
	}

	if len(tokens) == 0 {
		tokens = []int{0}
	}
	return tokens
}

func detokenize(tokens []neural.Token) string {
	if len(tokens) == 0 {
		return ""
	}

	mu.RLock()
	defer mu.RUnlock()

	words := make([]string, 0, len(tokens))
	for _, t := range tokens {
		if word, exists := idToWord[t.ID]; exists {
			if word != "<PAD>" && word != "<UNK>" {
				words = append(words, word)
			}
		}
	}

	return strings.Join(words, " ")
}

func SaveVocab(path string) error {
	mu.RLock()
	defer mu.RUnlock()

	file, err := os.Create(path)
	if err != nil {
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

	return gob.NewEncoder(file).Encode(data)
}

func LoadVocab(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var data struct {
		WordToID map[string]int
		IDToWord map[int]string
		NextID   int
	}

	if err := gob.NewDecoder(file).Decode(&data); err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()

	wordToID = data.WordToID
	idToWord = data.IDToWord
	nextID = data.NextID

	return nil
}
