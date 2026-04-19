package model

import (
	"math"
	"strings"

	"chatbot/internal/config"
)

func NewSimpleTextMatcher(conversations map[string]string) *SimpleTextMatcher {
	return &SimpleTextMatcher{Conversations: conversations}
}

func normalize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.TrimRight(s, ".,!?;:")
	return s
}

func (m *SimpleTextMatcher) FindBestMatch(input string) (string, float64) {
	normalizedInput := normalize(input)
	inputWords := strings.Fields(normalizedInput)

	if len(inputWords) == 0 {
		return "", 0
	}

	bestMatch, bestScore := "", 0.0
	for question, answer := range m.Conversations {
		normalizedQuestion := normalize(question)
		questionWords := strings.Fields(normalizedQuestion)

		score := m.calculateMatchScore(inputWords, questionWords)
		lenDiff := math.Abs(float64(len(inputWords) - len(questionWords)))
		maxLen := math.Max(float64(len(inputWords)), float64(len(questionWords)))
		score = score * (1.0 - (lenDiff/maxLen)*config.LengthPenalty)

		if score > bestScore {
			bestScore, bestMatch = score, answer
		}
	}

	if bestScore < config.MinSimilarityScore {
		return "", 0
	}
	return bestMatch, bestScore
}

func (m *SimpleTextMatcher) calculateMatchScore(inputWords, questionWords []string) float64 {
	if len(inputWords) == 0 || len(questionWords) == 0 {
		return 0
	}

	totalScore := 0.0
	matchedInput := make([]bool, len(inputWords))
	matchedQuestion := make([]bool, len(questionWords))

	for i, w1 := range inputWords {
		for j, w2 := range questionWords {
			if !matchedInput[i] && !matchedQuestion[j] && w1 == w2 {
				totalScore += 1.0
				matchedInput[i], matchedQuestion[j] = true, true
				break
			}
		}
	}

	for i, w1 := range inputWords {
		if matchedInput[i] {
			continue
		}
		for j, w2 := range questionWords {
			if !matchedQuestion[j] && (strings.Contains(w1, w2) || strings.Contains(w2, w1)) {
				totalScore += 0.5
				matchedQuestion[j] = true
				break
			}
		}
	}

	score := totalScore / float64(len(questionWords))
	if score > 1.0 {
		score = 1.0
	}
	return score
}

func (m *SimpleTextMatcher) ResolveSynonyms(input string, synonyms map[string]string) string {
	normalizedInput := normalize(input)

	if mainPhrase, exists := synonyms[normalizedInput]; exists {
		return mainPhrase
	}
	for _, word := range strings.Fields(normalizedInput) {
		if mainPhrase, exists := synonyms[word]; exists {
			return mainPhrase
		}
	}
	return input
}
