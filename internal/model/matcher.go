package model

import (
	"math"
	"strings"

	"chatbot/internal/config"
)

func (m *SimpleTextMatcher) FindBestMatch(input string) (string, float64) {
	inputWords := strings.Fields(normalize(input))
	if len(inputWords) == 0 {
		return "", 0
	}

	bestMatch, bestScore := "", 0.0
	for question, answer := range m.Conversations {
		questionWords := strings.Fields(normalize(question))
		score := m.calculateMatchScore(inputWords, questionWords)
		score = m.applyLengthPenalty(score, inputWords, questionWords)

		if score > bestScore {
			bestScore, bestMatch = score, answer
		}
	}

	if bestScore < config.MinSimilarityScore {
		return "", 0
	}
	return bestMatch, bestScore * 100
}

func (m *SimpleTextMatcher) calculateMatchScore(inputWords, questionWords []string) float64 {
	if len(inputWords) == 0 || len(questionWords) == 0 {
		return 0
	}

	totalScore := 0.0
	matchedQuestion := make([]bool, len(questionWords))

	for i, w1 := range inputWords {
		for j, w2 := range questionWords {
			if !matchedQuestion[j] && strings.EqualFold(w1, w2) {
				posBonus := 1.0 - (math.Abs(float64(i-j))/float64(max(len(inputWords), len(questionWords))))*0.5
				totalScore += 1.0 * posBonus
				matchedQuestion[j] = true
				break
			}
		}
	}

	return math.Min(1.0, totalScore/float64(len(questionWords)))
}

func (m *SimpleTextMatcher) applyLengthPenalty(score float64, inputWords, questionWords []string) float64 {
	lenDiff := math.Abs(float64(len(inputWords) - len(questionWords)))
	maxLen := math.Max(float64(len(inputWords)), float64(len(questionWords)))
	return score * (1.0 - (lenDiff/maxLen)*config.LengthPenalty)
}

func (m *SimpleTextMatcher) ResolveSynonyms(input string, synonyms map[string]string) string {
	normalized := normalize(input)

	if main, exists := synonyms[normalized]; exists {
		return main
	}
	for _, word := range strings.Fields(normalized) {
		if main, exists := synonyms[word]; exists {
			return main
		}
	}
	return input
}

func normalize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	return strings.TrimRight(s, ".,!?;:")
}
