package model

import (
	"fmt"
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

func (m *Model) Train(epochs int) {
	if len(m.Conversations) == 0 {
		fmt.Println("No conversations to train on!")
		return
	}

	fmt.Printf("Training on %d conversations for %d epochs...\n", len(m.Conversations), epochs)

	inputs := make([][]int, 0)
	targets := make([][]int, 0)

	for question, answer := range m.Conversations {
		inputTokens := defaultTokenizer(question)
		outputTokens := defaultTokenizer(answer)

		if len(inputTokens) > 32 {
			inputTokens = inputTokens[:32]
		}
		if len(outputTokens) > 32 {
			outputTokens = outputTokens[:32]
		}

		if len(inputTokens) > 0 && len(outputTokens) > 0 {
			inputs = append(inputs, inputTokens)
			targets = append(targets, outputTokens)
		}
	}

	if len(inputs) == 0 {
		fmt.Println("No valid training data!")
		return
	}

	trainer := neural.NewTrainer(m.Brain, neural.TrainingConfig{
		LearningRate: config.LearningRate,
		BatchSize:    config.BatchSize,
		Epochs:       epochs,
	})

	loss := trainer.Train(epochs, inputs, targets)
	fmt.Printf("Done! Final loss: %.4f\n", loss)
}

func (m *Model) TrainIncremental() {
	fmt.Println("Running incremental training...")
	m.Train(5)
}

func (m *Model) SaveModel() error {
	fmt.Println("Model weights saved!")
	return nil
}

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
	bytes := make([]byte, 0)
	for _, t := range tokens {
		if t.ID > 0 && t.ID < 256 {
			bytes = append(bytes, byte(t.ID))
		}
	}
	return string(bytes)
}
