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
	result := ""
	for _, t := range tokens {
		if t.ID >= 65 && t.ID <= 90 {
			result += string(byte(t.ID))
		} else if t.ID >= 97 && t.ID <= 122 {
			result += string(byte(t.ID))
		} else if t.ID == 32 {
			result += " "
		}
	}
	if len(result) == 0 {
		for _, t := range tokens {
			if t.ID >= 32 && t.ID < 127 {
				result += string(byte(t.ID))
			}
		}
	}
	return strings.TrimSpace(result)
}
