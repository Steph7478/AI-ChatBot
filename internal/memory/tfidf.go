package memory

import (
	"math"

	"chatbot/internal/config"
)

func (m *ConversationMemory) CalculateIDF() {
	docCount := make(map[int]int)

	for _, tf := range m.QuestionTF {
		for token := range tf {
			docCount[token]++
		}
	}

	numDocs := len(m.Questions)
	for token, count := range docCount {
		m.IDF[token] = math.Log(float64(numDocs) / (float64(count) + config.IDFSmoothing))
	}
}

func (m *ConversationMemory) GetTFIDFVector(tokens []int) []float64 {
	vector := make([]float64, m.VocabularySize)

	tf := m.calculateTF(tokens)

	for token, freq := range tf {
		idf := m.IDF[token]
		if idf == 0 {
			idf = 1.0
		}
		if token >= 0 && token < len(vector) {
			vector[token] = freq * idf
		}
	}

	return vector
}

func (m *ConversationMemory) calculateTF(tokens []int) map[int]float64 {
	tf := make(map[int]float64)
	for _, token := range tokens {
		tf[token]++
	}

	maxFreq := 0.0
	for _, count := range tf {
		if count > maxFreq {
			maxFreq = count
		}
	}

	if maxFreq > 0 {
		for token := range tf {
			tf[token] = tf[token] / maxFreq
		}
	}

	return tf
}

func (m *ConversationMemory) CosineSimilarity(vec1, vec2 []float64) float64 {
	if len(vec1) != len(vec2) {
		return 0
	}

	dot := 0.0
	norm1 := 0.0
	norm2 := 0.0

	for i := range vec1 {
		dot += vec1[i] * vec2[i]
		norm1 += vec1[i] * vec1[i]
		norm2 += vec2[i] * vec2[i]
	}

	if norm1 == 0 || norm2 == 0 {
		return 0
	}

	return dot / (math.Sqrt(norm1) * math.Sqrt(norm2))
}
