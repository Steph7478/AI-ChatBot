package memory

import (
	"encoding/gob"
	"math"
	"math/rand"
	"os"
	"sort"

	"chatbot/internal/config"
	"chatbot/internal/dataset"
)

type ConversationMemory struct {
	Questions      [][]int
	Answers        [][]int
	QuestionText   []string
	AnswerText     []string
	Vocab          *dataset.Vocabulary
	QuestionTF     []map[int]float64
	IDF            map[int]float64
	VocabularySize int
}

func NewConversationMemory(vocab *dataset.Vocabulary) *ConversationMemory {
	return &ConversationMemory{
		Questions:      make([][]int, 0),
		Answers:        make([][]int, 0),
		QuestionText:   make([]string, 0),
		AnswerText:     make([]string, 0),
		Vocab:          vocab,
		QuestionTF:     make([]map[int]float64, 0),
		IDF:            make(map[int]float64),
		VocabularySize: vocab.Size,
	}
}

func (m *ConversationMemory) Learn(questionTokens, answerTokens []int) {
	m.Questions = append(m.Questions, questionTokens)
	m.Answers = append(m.Answers, answerTokens)

	m.QuestionText = append(m.QuestionText, m.Vocab.Detokenize(questionTokens))
	m.AnswerText = append(m.AnswerText, m.Vocab.Detokenize(answerTokens))

	tf := make(map[int]float64)
	for _, token := range questionTokens {
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
	m.QuestionTF = append(m.QuestionTF, tf)
}

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

	type candidate struct {
		index int
		score float64
	}

	candidates := make([]candidate, 0)

	for i := 0; i < len(m.Questions); i++ {
		questionVector := m.GetTFIDFVector(m.Questions[i])
		tfidfScore := m.CosineSimilarity(inputVector, questionVector)

		finalScore := tfidfScore
		if i == exactIdx && exactScore > config.MinWordMatchScore {
			finalScore = exactScore * config.ExactMatchBoost
		}

		if finalScore > config.SimilarityThreshold {
			candidates = append(candidates, candidate{index: i, score: finalScore})
		}
	}

	if len(candidates) == 0 {
		return nil
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].score > candidates[j].score
	})

	bestScore := candidates[0].score

	if bestScore < config.SimilarityThreshold {
		return nil
	}

	topK := config.TopKCandidates
	if temperature > 0.7 {
		topK = config.TopKTempHigh
	} else if temperature < 0.3 {
		topK = config.TopKTempLow
	}

	filtered := make([]candidate, 0)
	for _, c := range candidates {
		if c.score >= bestScore*0.8 && len(filtered) < topK {
			filtered = append(filtered, c)
		}
	}

	if len(filtered) == 0 {
		filtered = candidates[:min(topK, len(candidates))]
	}

	var selected candidate
	if temperature > 0.5 && len(filtered) > 1 {
		probs := make([]float64, len(filtered))
		sum := 0.0
		for i, cand := range filtered {
			probs[i] = math.Exp(cand.score / temperature)
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
				selected = filtered[i]
				break
			}
		}
	} else {
		selected = filtered[0]
	}

	return m.Answers[selected.index]
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

func (m *ConversationMemory) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(m)
}

func (m *ConversationMemory) Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	return decoder.Decode(m)
}
