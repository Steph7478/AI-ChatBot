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
		fmt.Printf("\n🤖 Bot [EXACT MATCH]: %s\n", response)
		return ResponseResult{
			Text:       response,
			Type:       ResponseGenerated,
			Confidence: 0.95,
		}
	}

	if match, score := m.Matcher.FindBestMatch(prompt); match != "" && score > config.MinSimilarityScore {
		fmt.Printf("\n🤖 Bot [FUZZY MATCH - %.1f%%]: %s\n", score*100, match)
		return ResponseResult{
			Text:       match,
			Type:       ResponseGenerated,
			Confidence: score,
		}
	}

	response, isFallback := m.generateFromNeural(prompt, temp)

	if isFallback {
		fmt.Printf("\n🤖 Bot [NEURAL - FALLBACK]: %s\n", response)
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

	fmt.Printf("\n🤖 Bot [NEURAL GENERATED]: %s\n", response)
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

	for _, t := range resp.Tokens {
		fmt.Printf("%d ", t.ID)
	}
	fmt.Println()

	response := detokenize(resp.Tokens)

	words := strings.Fields(response)
	if len(words) > 2 {
		filtered := []string{words[0]}
		for i := 1; i < len(words); i++ {
			if words[i] != filtered[len(filtered)-1] {
				filtered = append(filtered, words[i])
			}
		}
		response = strings.Join(filtered, " ")
	}

	if len(response) < 3 {
		return "I need more training to answer this properly.", true
	}

	fmt.Println("[DEBUG] Returning generated response")
	return response, false
}
