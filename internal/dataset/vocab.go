package dataset

import (
	"strings"
	"unicode"
)

func NewVocabulary() *Vocabulary {
	v := &Vocabulary{
		WordToID: make(map[string]int),
		IDToWord: make(map[int]string),
	}
	v.AddWord("<pad>")
	v.AddWord("<unk>")
	v.AddWord("<eos>")
	return v
}

func (v *Vocabulary) AddWord(word string) int {
	if id, ok := v.WordToID[word]; ok {
		return id
	}
	id := v.Size
	v.WordToID[word] = id
	v.IDToWord[id] = word
	v.Size++
	return id
}

func (v *Vocabulary) GetID(word string) int {
	if id, ok := v.WordToID[word]; ok {
		return id
	}
	return v.WordToID["<unk>"]
}

func (v *Vocabulary) Tokenize(text string) []int {
	words := tokenize(text)
	tokens := make([]int, len(words))
	for i, w := range words {
		tokens[i] = v.GetID(w)
	}
	return tokens
}

func (v *Vocabulary) AddAndTokenize(text string) []int {
	words := tokenize(text)
	tokens := make([]int, len(words))
	for i, w := range words {
		tokens[i] = v.AddWord(w)
	}
	return tokens
}

func (v *Vocabulary) Detokenize(tokens []int) string {
	words := make([]string, 0, len(tokens))
	for _, id := range tokens {
		if word, ok := v.IDToWord[id]; ok && !strings.HasPrefix(word, "<") {
			words = append(words, word)
		}
	}
	return strings.Join(words, " ")
}

func tokenize(text string) []string {
	text = strings.ToLower(text)
	var b strings.Builder
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' {
			b.WriteRune(r)
		} else {
			b.WriteRune(' ')
		}
	}
	return strings.Fields(b.String())
}
