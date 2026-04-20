package model

import (
	"bufio"
	"fmt"
	"strings"

	"chatbot/internal/config"
	"chatbot/internal/core"
)

const (
	cyan  = "\033[36m"
	gray  = "\033[90m"
	reset = "\033[0m"
)

func (m *Model) GenerateResponse(prompt string, temp float64, scanner *bufio.Scanner) ResponseResult {
	originalPrompt := prompt

	if response, exists := m.Conversations[originalPrompt]; exists {
		m.printResponse(response, "exact match")
		return ResponseResult{Text: response, Type: ResponseGenerated, Confidence: 0.95}
	}

	if match, score := m.Matcher.FindBestMatch(originalPrompt); match != "" && score > config.MinSimilarityScore*100 {
		m.printResponse(match, fmt.Sprintf("fuzzy match - %.1f%%", score))
		return ResponseResult{Text: match, Type: ResponseGenerated, Confidence: score / 100}
	}

	resolved := m.Matcher.ResolveSynonyms(originalPrompt, m.Synonyms)
	if resolved != originalPrompt {
		if response, exists := m.Conversations[resolved]; exists {
			m.printResponse(response, "synonym match")
			return ResponseResult{Text: response, Type: ResponseGenerated, Confidence: 0.95}
		}
		if match, score := m.Matcher.FindBestMatch(resolved); match != "" && score > config.MinSimilarityScore*100 {
			m.printResponse(match, fmt.Sprintf("synonym fuzzy - %.1f%%", score))
			return ResponseResult{Text: match, Type: ResponseGenerated, Confidence: score / 100}
		}
	}

	response, isFallback := m.generateFromNeural(originalPrompt, temp)

	if isFallback {
		m.handleFallback(originalPrompt, response, scanner)
		return ResponseResult{Text: response, Type: ResponseFallback, Confidence: 0}
	}

	m.handleNeuralResponse(originalPrompt, response, scanner)
	return ResponseResult{Text: response, Type: ResponseGenerated, Confidence: 0.8}
}

func (m *Model) printResponse(response, matchType string) {
	fmt.Printf("%sBot: %s%s\n", cyan, response, reset)
	fmt.Printf("%s[%s]%s\n", gray, matchType, reset)
}

func (m *Model) generateFromNeural(prompt string, temp float64) (string, bool) {
	cfg := core.InferenceConfig{
		Temperature: temp,
		TopK:        config.TopK,
		MaxTokens:   config.MaxTokens,
	}

	resp := m.Brain.Generate(prompt, DefaultTokenizer, cfg)
	response := Detokenize(resp.Tokens)

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

func (m *Model) handleFallback(prompt, response string, scanner *bufio.Scanner) {
	fmt.Printf("%sBot: %s%s\n", cyan, response, reset)
	fmt.Printf("%s[fallback - please teach me]%s\n", gray, reset)
	fmt.Print("📝 Please teach me the correct response: ")
	scanner.Scan()
	correct := strings.TrimSpace(scanner.Text())
	if correct != "" {
		m.LearnForTraining(prompt, correct)
		fmt.Println("✅ Thanks for teaching me! 🧠")
	} else {
		fmt.Println("❌ No response saved.")
	}
}

func (m *Model) handleNeuralResponse(prompt, response string, scanner *bufio.Scanner) {
	fmt.Printf("%sBot: %s%s\n", cyan, response, reset)
	fmt.Printf("%s[generated via neural]%s\n", gray, reset)
	fmt.Print("❓ Was this response useful? (y/n): ")
	scanner.Scan()
	answer := strings.ToLower(strings.TrimSpace(scanner.Text()))

	if answer == "y" || answer == "yes" {
		m.LearnForTraining(prompt, response)
		fmt.Println("✅ Thanks! I'll remember that! 🧠")
	} else {
		fmt.Print("📝 What would be the correct response? ")
		scanner.Scan()
		correct := strings.TrimSpace(scanner.Text())
		if correct != "" {
			m.LearnForTraining(prompt, correct)
			fmt.Println("✅ Thanks for teaching me! 🧠")
		} else {
			fmt.Println("❌ No response saved.")
		}
	}
}

func (m *Model) LearnForTraining(input, response string) {
	m.TrainingData[input] = response
	m.SaveTrainingData(input, response)
}
