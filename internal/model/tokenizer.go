// internal/model/tokenizer.go
package model

import (
	"strings"

	"chatbot/internal/config"
	"chatbot/internal/neural"
)

func defaultTokenizer(text string) []int {
	if text == "" {
		return []int{0}
	}
	tokens := make([]int, 0)
	for _, r := range text {
		token := int(r) % config.VocabSize
		tokens = append(tokens, token)
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
	result := make([]byte, 0)
	for _, t := range tokens {
		if t.ID >= 32 && t.ID <= 126 {
			result = append(result, byte(t.ID))
		}
	}
	return strings.TrimSpace(string(result))
}
