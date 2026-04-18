package model

import (
	"strings"

	"chatbot/internal/config"
	"chatbot/internal/neural"
)

type Model struct {
	Brain         *neural.Transformer
	Matcher       *SimpleTextMatcher
	Conversations map[string]string
	Synonyms      map[string]string
}

func NewModel() *Model {
	m := &Model{
		Brain: neural.NewTransformer(
			config.VocabSize,
			config.EmbeddingDim,
			config.HiddenDim,
			config.NumHeads,
			config.NumLayers,
			config.MaxSeqLen,
			config.Dropout,
		),
		Conversations: make(map[string]string),
		Synonyms:      make(map[string]string),
	}

	m.LoadAll()
	m.Matcher = NewSimpleTextMatcher(m.Conversations)

	return m
}

func (m *Model) normalize(text string) string {
	punctuation := []string{"?", "!", ".", ",", ";", ":", "'", "\"", "(", ")"}
	result := strings.ToLower(strings.TrimSpace(text))

	for _, p := range punctuation {
		result = strings.ReplaceAll(result, p, "")
	}

	result = strings.Join(strings.Fields(result), " ")

	if main, ok := m.Synonyms[result]; ok {
		return main
	}
	return result
}

func (m *Model) Learn(input, response string) {
	key := m.normalize(input)
	m.Conversations[key] = response
	m.SaveConversation(input, response)
}

func (m *Model) GenerateResponse(prompt string, temp float64) ResponseResult {
	normalized := m.normalize(prompt)

	if response, exists := m.Conversations[normalized]; exists {
		return ResponseResult{
			Text:       response,
			Type:       ResponseGenerated,
			Confidence: 0.95,
		}
	}

	if match, score := m.Matcher.FindBestMatch(prompt); match != "" && score > 0.4 {
		return ResponseResult{
			Text:       match,
			Type:       ResponseGenerated,
			Confidence: score,
		}
	}

	cfg := neural.InferenceConfig{
		Temperature: temp,
		TopK:        config.TopK,
		MaxTokens:   config.MaxTokens,
	}

	resp := m.Brain.Generate(prompt, defaultTokenizer, cfg)

	if resp.Confidence < config.MinConfidence {
		return ResponseResult{
			Text:       "I'm still learning. Please teach me! ☕",
			Type:       ResponseFallback,
			Confidence: resp.Confidence,
		}
	}

	response := detokenize(resp.Tokens)
	if len(response) < 2 {
		return ResponseResult{
			Text:       "I'm still learning. Please teach me! ☕",
			Type:       ResponseFallback,
			Confidence: 0,
		}
	}

	return ResponseResult{
		Text:       response,
		Type:       ResponseGenerated,
		Confidence: resp.Confidence,
	}
}

func defaultTokenizer(text string) []int {
	tokens := make([]int, len(text))
	for i, r := range text {
		tokens[i] = int(r) % config.VocabSize
	}
	return tokens
}

func detokenize(tokens []neural.Token) string {
	bytes := make([]byte, len(tokens))
	for i, t := range tokens {
		bytes[i] = byte(t.ID % 256)
	}
	return string(bytes)
}
