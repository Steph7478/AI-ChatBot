package model

import (
	"strings"
)

type SimpleTextMatcher struct {
	Conversations map[string]string
}

func NewSimpleTextMatcher(conversations map[string]string) *SimpleTextMatcher {
	return &SimpleTextMatcher{
		Conversations: conversations,
	}
}

func (m *SimpleTextMatcher) FindBestMatch(input string) (string, float64) {
	input = strings.ToLower(strings.TrimSpace(input))
	inputWords := strings.Fields(input)

	bestMatch := ""
	bestScore := 0.0

	for question, answer := range m.Conversations {
		questionWords := strings.Fields(question)

		common := 0
		for _, w1 := range inputWords {
			for _, w2 := range questionWords {
				if w1 == w2 || strings.Contains(w1, w2) || strings.Contains(w2, w1) {
					common++
					break
				}
			}
		}

		totalUnique := len(unique(append(inputWords, questionWords...)))
		if totalUnique > 0 {
			score := float64(common) / float64(totalUnique)
			if score > bestScore {
				bestScore = score
				bestMatch = answer
			}
		}
	}

	if bestScore < 0.3 {
		return "", 0
	}

	return bestMatch, bestScore
}

func unique(words []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, w := range words {
		if !seen[w] {
			seen[w] = true
			result = append(result, w)
		}
	}
	return result
}
