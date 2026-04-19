package model

import (
	"bufio"
	"fmt"
	"strings"

	"chatbot/internal/config"
	"chatbot/internal/neural"
)

const (
	cyan  = "\033[36m"
	gray  = "\033[90m"
	reset = "\033[0m"
)

func (m *Model) GenerateResponse(prompt string, temp float64, scanner *bufio.Scanner) ResponseResult {
	originalPrompt := prompt

	if response, exists := m.Conversations[originalPrompt]; exists {
		fmt.Printf("%sBot: %s%s\n", cyan, response, reset)
		fmt.Printf("%s[exact match]%s\n", gray, reset)
		return ResponseResult{
			Text:       response,
			Type:       ResponseGenerated,
			Confidence: 0.95,
		}
	}

	if match, score := m.Matcher.FindBestMatch(originalPrompt); match != "" && score > config.MinSimilarityScore {
		fmt.Printf("%sBot: %s%s\n", cyan, match, reset)
		fmt.Printf("%s[fuzzy match - %.1f%%]%s\n", gray, score*100, reset)
		return ResponseResult{
			Text:       match,
			Type:       ResponseGenerated,
			Confidence: score,
		}
	}

	resolved := m.Matcher.ResolveSynonyms(originalPrompt, m.Synonyms)
	if resolved != originalPrompt {
		if response, exists := m.Conversations[resolved]; exists {
			fmt.Printf("%sBot: %s%s\n", cyan, response, reset)
			fmt.Printf("%s[synonym match]%s\n", gray, reset)
			return ResponseResult{
				Text:       response,
				Type:       ResponseGenerated,
				Confidence: 0.95,
			}
		}
		if match, score := m.Matcher.FindBestMatch(resolved); match != "" && score > config.MinSimilarityScore {
			fmt.Printf("%sBot: %s%s\n", cyan, match, reset)
			fmt.Printf("%s[synonym fuzzy - %.1f%%]%s\n", gray, score*100, reset)
			return ResponseResult{
				Text:       match,
				Type:       ResponseGenerated,
				Confidence: score,
			}
		}
	}

	response, isFallback := m.generateFromNeural(originalPrompt, temp)

	if isFallback {
		fmt.Printf("%sBot: %s%s\n", cyan, response, reset)
		fmt.Printf("%s[fallback - please teach me]%s\n", gray, reset)
		fmt.Print("📝 Please teach me the correct response: ")
		scanner.Scan()
		correct := strings.TrimSpace(scanner.Text())
		if correct != "" {
			m.Learn(originalPrompt, correct)
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

	fmt.Printf("%sBot: %s%s\n", cyan, response, reset)
	fmt.Printf("%s[generated via neural]%s\n", gray, reset)
	fmt.Print("❓ Was this response useful? (y/n): ")
	scanner.Scan()
	answer := strings.ToLower(strings.TrimSpace(scanner.Text()))

	if answer == "y" || answer == "yes" {
		m.Learn(originalPrompt, response)
		fmt.Println("✅ Thanks! I'll remember that! 🧠")
	} else {
		fmt.Print("📝 What would be the correct response? ")
		scanner.Scan()
		correct := strings.TrimSpace(scanner.Text())
		if correct != "" {
			m.Learn(originalPrompt, correct)
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

	return response, false
}

func (m *Model) Learn(input, response string) {
	m.Conversations[input] = response
	m.SaveConversation(input, response)
}

func (m *Model) LearnAndSave(input, response string) error {
	m.Learn(input, response)
	return m.SaveConversation(input, response)
}
