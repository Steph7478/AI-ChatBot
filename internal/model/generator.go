package model

import (
	"bufio"
	"fmt"
	"strings"

	"chatbot/internal/config"
	"chatbot/internal/neural"
)

func (m *Model) GenerateResponse(prompt string, temp float64, scanner *bufio.Scanner) ResponseResult {
	if response, exists := m.Conversations[prompt]; exists {
		fmt.Printf("\n🤖 Bot: %s\n", response)
		return ResponseResult{
			Text:       response,
			Type:       ResponseGenerated,
			Confidence: 0.95,
		}
	}

	if match, score := m.Matcher.FindBestMatch(prompt); match != "" && score > config.MinSimilarityScore {
		fmt.Printf("\n🤖 Bot: %s\n", match)
		return ResponseResult{
			Text:       match,
			Type:       ResponseGenerated,
			Confidence: score,
		}
	}

	response, isFallback := m.generateFromNeural(prompt, temp)

	fmt.Printf("\n🤖 Bot: %s\n", response)

	if isFallback {
		fmt.Print("📝 Please teach me the correct response: ")
		scanner.Scan()
		correct := strings.TrimSpace(scanner.Text())
		if correct != "" {
			m.Learn(prompt, correct)
			fmt.Println("✅ Thanks for teaching me! 🧠")
		} else {
			fmt.Println("❌ No response saved.")
		}
		return ResponseResult{
			Text:       response,
			Type:       ResponseFallback,
			Confidence: 0,
		}
	}

	fmt.Print("❓ Was this response useful? (y/n): ")
	scanner.Scan()
	answer := strings.ToLower(strings.TrimSpace(scanner.Text()))

	if answer == "y" || answer == "yes" {
		m.Learn(prompt, response)
		fmt.Println("✅ Thanks! I'll remember that! 🧠")
	} else {
		fmt.Print("📝 What would be the correct response? ")
		scanner.Scan()
		correct := strings.TrimSpace(scanner.Text())
		if correct != "" {
			m.Learn(prompt, correct)
			fmt.Println("✅ Thanks for teaching me! 🧠")
		} else {
			fmt.Println("❌ No response saved.")
		}
	}

	return ResponseResult{
		Text:       response,
		Type:       ResponseGenerated,
		Confidence: 0.8,
	}
}

func (m *Model) generateFromNeural(prompt string, temp float64) (string, bool) {
	cfg := neural.InferenceConfig{
		Temperature: temp,
		TopK:        config.TopK,
		MaxTokens:   config.MaxTokens,
	}

	resp := m.Brain.Generate(prompt, defaultTokenizer, cfg)
	response := detokenize(resp.Tokens)

	temperatures := []float64{1.2, 1.5, 1.8}
	topKs := []int{40, 60, 80}

	for i := 0; i < len(temperatures) && len(response) < 3; i++ {
		cfg.Temperature = temperatures[i]
		cfg.TopK = topKs[i]
		resp = m.Brain.Generate(prompt, defaultTokenizer, cfg)
		response = detokenize(resp.Tokens)
	}

	if len(response) < 3 {
		return "I need more training to answer this properly.", true
	}

	return response, false
}
