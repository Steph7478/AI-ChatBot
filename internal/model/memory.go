package model

import (
	"math"
	"math/rand"
	"sort"

	"chatbot/internal/config"
	"chatbot/internal/dataset"
)

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

func (m *ConversationMemory) Learn(q, a []int) {
	m.Questions = append(m.Questions, q)
	m.Answers = append(m.Answers, a)
	m.QuestionText = append(m.QuestionText, m.Vocab.Detokenize(q))
	m.AnswerText = append(m.AnswerText, m.Vocab.Detokenize(a))
	m.QuestionTF = append(m.QuestionTF, m.tf(q))
}

func (m *ConversationMemory) tf(tokens []int) map[int]float64 {
	tf := make(map[int]float64)
	for _, t := range tokens {
		tf[t]++
	}
	var max float64
	for _, c := range tf {
		if c > max {
			max = c
		}
	}
	if max > 0 {
		for t := range tf {
			tf[t] /= max
		}
	}
	return tf
}

func (m *ConversationMemory) CalculateIDF() {
	docCnt := make(map[int]int)
	for _, tf := range m.QuestionTF {
		for t := range tf {
			docCnt[t]++
		}
	}
	n := len(m.Questions)
	for t, c := range docCnt {
		m.IDF[t] = math.Log(float64(n) / (float64(c) + config.IDFSmoothing))
	}
}

func (m *ConversationMemory) vector(tokens []int) []float64 {
	vec := make([]float64, m.VocabularySize)
	tf := m.tf(tokens)
	for t, f := range tf {
		idf := m.IDF[t]
		if idf == 0 {
			idf = 1
		}
		if t >= 0 && t < len(vec) {
			vec[t] = f * idf
		}
	}
	return vec
}

func (m *ConversationMemory) CosineSimilarity(v1, v2 []float64) float64 {
	if len(v1) != len(v2) {
		return 0
	}
	var dot, n1, n2 float64
	for i := range v1 {
		dot += v1[i] * v2[i]
		n1 += v1[i] * v1[i]
		n2 += v2[i] * v2[i]
	}
	if n1 == 0 || n2 == 0 {
		return 0
	}
	return dot / (math.Sqrt(n1) * math.Sqrt(n2))
}

func (m *ConversationMemory) exactMatch(q []int) (int, float64) {
	bestIdx, bestScore := -1, 0.0
	for i, quest := range m.Questions {
		seen := make(map[int]bool)
		common := 0
		for _, t := range q {
			for _, qt := range quest {
				if t == qt && !seen[t] {
					common++
					seen[t] = true
					break
				}
			}
		}
		if len(q) > 0 {
			score := float64(common) / float64(len(q))
			if score > bestScore {
				bestScore, bestIdx = score, i
			}
		}
	}
	return bestIdx, bestScore
}

func (m *ConversationMemory) FindResponse(q []int, temp float64) []int {
	if len(m.Questions) == 0 {
		return nil
	}
	if len(m.IDF) == 0 {
		m.CalculateIDF()
	}

	exactIdx, exactScore := m.exactMatch(q)
	inputVec := m.vector(q)

	type cand struct {
		idx   int
		score float64
	}
	var cands []cand

	for i := range m.Questions {
		score := m.CosineSimilarity(inputVec, m.vector(m.Questions[i]))
		if i == exactIdx && exactScore > config.MinWordMatchScore {
			score = exactScore * config.ExactMatchBoost
		}
		if score > config.SimilarityThreshold {
			cands = append(cands, cand{i, score})
		}
	}
	if len(cands) == 0 {
		return nil
	}

	sort.Slice(cands, func(i, j int) bool { return cands[i].score > cands[j].score })

	topK := config.TopKCandidates
	if temp > 0.7 {
		topK = config.TopKTempHigh
	} else if temp < 0.3 {
		topK = config.TopKTempLow
	}

	var filtered []cand
	best := cands[0].score
	for _, c := range cands {
		if c.score >= best*0.8 && len(filtered) < topK {
			filtered = append(filtered, c)
		}
	}
	if len(filtered) == 0 {
		filtered = cands[:min(topK, len(cands))]
	}

	var selected cand
	if temp > 0.5 && len(filtered) > 1 {
		var probs []float64
		var sum float64
		for _, c := range filtered {
			p := math.Exp(c.score / temp)
			probs = append(probs, p)
			sum += p
		}
		for i := range probs {
			probs[i] /= sum
		}
		r := rand.Float64()
		var cum float64
		for i, p := range probs {
			cum += p
			if r < cum {
				selected = filtered[i]
				break
			}
		}
	} else {
		selected = filtered[0]
	}
	return m.Answers[selected.idx]
}
