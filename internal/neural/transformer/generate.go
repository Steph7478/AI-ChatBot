package transformer

import (
	"chatbot/internal/core"
)

func (t *Transformer) Generate(input string, tokenizer func(string) []int, cfg core.InferenceConfig) core.Response {
	ids := tokenizer(input)
	if len(ids) == 0 {
		return core.Response{}
	}

	tokens := []core.Token{}
	current := make([]int, len(ids))
	copy(current, ids)
	recent := make(map[int]int)

	maxTokens := cfg.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 20
	}

	for step := 0; step < maxTokens; step++ {
		logits := t.Forward(current)
		if len(logits) == 0 || len(logits[len(logits)-1]) == 0 {
			break
		}

		last := make([]float64, len(logits[len(logits)-1]))
		copy(last, logits[len(logits)-1])

		for i := range last {
			last[i] /= cfg.Temperature
		}

		probs := core.Softmax(last)

		for i := range probs {
			if count, ok := recent[i]; ok && count > 0 {
				probs[i] = probs[i] / float64(1+count)
			}
		}

		id := t.sample(probs, cfg.TopK)
		tokens = append(tokens, core.Token{ID: id, Prob: probs[id]})
		recent[id] = recent[id] + 1
		current = append(current, id)

		if len(current) > t.MaxSeqLen {
			current = current[1:]
		}

		if id == 0 {
			break
		}
	}

	if len(tokens) == 0 {
		return core.Response{}
	}

	confidence := tokens[len(tokens)-1].Prob
	for _, t := range tokens {
		if t.Prob < confidence {
			confidence = t.Prob
		}
	}

	return core.Response{Tokens: tokens, Confidence: confidence}
}

func (t *Transformer) sample(probs []float64, k int) int {
	if k <= 0 || k >= len(probs) {
		maxProb := 0.0
		id := 0
		for i, p := range probs {
			if p > maxProb {
				maxProb, id = p, i
			}
		}
		return id
	}

	indices := make([]int, len(probs))
	for i := range indices {
		indices[i] = i
	}
	for i := 0; i < k; i++ {
		maxIdx := i
		for j := i + 1; j < len(probs); j++ {
			if probs[indices[j]] > probs[indices[maxIdx]] {
				maxIdx = j
			}
		}
		indices[i], indices[maxIdx] = indices[maxIdx], indices[i]
	}
	return indices[0]
}
