package dataset

import (
	"strings"
)

type Vocabulary struct {
	WordToID map[string]int
	IDToWord map[int]string
	Size     int
}

func NewVocabulary() *Vocabulary {
	vocab := &Vocabulary{
		WordToID: make(map[string]int),
		IDToWord: make(map[int]string),
	}
	vocab.AddWord("<pad>")
	vocab.AddWord("<unk>")
	vocab.AddWord("<eos>")
	return vocab
}

func (v *Vocabulary) AddWord(word string) int {
	if id, exists := v.WordToID[word]; exists {
		return id
	}
	id := v.Size
	v.WordToID[word] = id
	v.IDToWord[id] = word
	v.Size++
	return id
}

func (v *Vocabulary) GetID(word string) int {
	if id, exists := v.WordToID[word]; exists {
		return id
	}
	return v.WordToID["<unk>"]
}

func (v *Vocabulary) Detokenize(tokens []int) string {
	words := make([]string, 0, len(tokens))
	for _, id := range tokens {
		if word, exists := v.IDToWord[id]; exists {
			if !strings.HasPrefix(word, "<") {
				words = append(words, word)
			}
		}
	}
	return strings.Join(words, " ")
}
