package memory

import (
	"math"
	"math/rand"
	"sort"

	"chatbot/internal/config"
)

func (m *ConversationMemory) exactWordMatch(questionTokens []int) (int, float64) {
	bestIdx := -1
	bestScore := 0.0

	for i, q := range m.Questions {
		common := 0
		seen := make(map[int]bool)
		for _, token := range questionTokens {
			for _, qToken := range q {
				if token == qToken && !seen[token] {
					common++
					seen[token] = true
					break
				}
			}
		}

		if len(questionTokens) > 0 {
			score := float64(common) / float64(len(questionTokens))
			if score > bestScore {
				bestScore = score
				bestIdx = i
			}
		}
	}

	return bestIdx, bestScore
}

func (m *ConversationMemory) FindResponse(questionTokens []int, temperature float64) []int {
	if len(m.Questions) == 0 {
		return nil
	}

	if len(m.IDF) == 0 {
		m.CalculateIDF()
	}

	exactIdx, exactScore := m.exactWordMatch(questionTokens)
	inputVector := m.GetTFIDFVector(questionTokens)

	candidates := m.evaluateCandidates(inputVector, exactIdx, exactScore)

	if len(candidates) == 0 {
		return nil
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Score > candidates[j].Score
	})

	bestScore := candidates[0].Score

	if bestScore < config.SimilarityThreshold {
		return nil
	}

	filtered := m.filterTopCandidates(candidates, bestScore, temperature)

	if len(filtered) == 0 {
		filtered = candidates[:min(config.TopKCandidates, len(candidates))]
	}

	selected := m.selectCandidate(filtered, temperature)

	return m.Answers[selected.Index]
}

func (m *ConversationMemory) evaluateCandidates(inputVector []float64, exactIdx int, exactScore float64) []Candidate {
	candidates := make([]Candidate, 0)

	for i := 0; i < len(m.Questions); i++ {
		questionVector := m.GetTFIDFVector(m.Questions[i])
		tfidfScore := m.CosineSimilarity(inputVector, questionVector)

		finalScore := tfidfScore
		if i == exactIdx && exactScore > config.MinWordMatchScore {
			finalScore = exactScore * config.ExactMatchBoost
		}

		if finalScore > config.SimilarityThreshold {
			candidates = append(candidates, Candidate{Index: i, Score: finalScore})
		}
	}

	return candidates
}

func (m *ConversationMemory) filterTopCandidates(candidates []Candidate, bestScore float64, temperature float64) []Candidate {
	topK := config.TopKCandidates
	if temperature > 0.7 {
		topK = config.TopKTempHigh
	} else if temperature < 0.3 {
		topK = config.TopKTempLow
	}

	filtered := make([]Candidate, 0)
	for _, c := range candidates {
		if c.Score >= bestScore*0.8 && len(filtered) < topK {
			filtered = append(filtered, c)
		}
	}

	return filtered
}

func (m *ConversationMemory) selectCandidate(candidates []Candidate, temperature float64) Candidate {
	if temperature > 0.5 && len(candidates) > 1 {
		return m.weightedRandomSelect(candidates, temperature)
	}
	return candidates[0]
}

func (m *ConversationMemory) weightedRandomSelect(candidates []Candidate, temperature float64) Candidate {
	probs := make([]float64, len(candidates))
	sum := 0.0

	for i, cand := range candidates {
		probs[i] = math.Exp(cand.Score / temperature)
		sum += probs[i]
	}

	for i := range probs {
		probs[i] /= sum
	}

	r := rand.Float64()
	cumsum := 0.0
	for i, p := range probs {
		cumsum += p
		if r < cumsum {
			return candidates[i]
		}
	}

	return candidates[0]
}

func (m *ConversationMemory) Generate(questionTokens []int, maxTokens int, temperature float64) []int {
	response := m.FindResponse(questionTokens, temperature)
	if response == nil {
		return nil
	}

	if len(response) > maxTokens {
		return response[:maxTokens]
	}
	return response
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
