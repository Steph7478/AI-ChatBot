package model

import (
	"fmt"
	"strings"

	"chatbot/internal/config"
)

func NewLCSMatcher() *LCSMatcher {
	return &LCSMatcher{
		Documents: make([][]string, 0),
		Answers:   make([][]int, 0),
	}
}

func (l *LCSMatcher) tokenize(text string) []string {
	text = strings.ToLower(text)
	r := strings.NewReplacer(
		".", " ", "!", " ", "?", " ",
		",", " ", ";", " ", ":", " ",
		"\"", " ", "'", " ", "(", " ",
		")", " ", "[", " ", "]", " ",
		"{", " ", "}", " ",
	)
	return strings.Fields(r.Replace(text))
}

func (l *LCSMatcher) lcsLength(a, b []string) int {
	m, n := len(a), len(b)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}
	return dp[m][n]
}

func (l *LCSMatcher) lcsSimilarity(a, b []string) float64 {
	if len(a) == 0 {
		return 0
	}

	lcsLen := l.lcsLength(a, b)

	basePercent := (float64(lcsLen) / float64(len(a))) * 100
	diffPenalty := 0.0
	for _, word := range a {
		found := false
		for _, w := range b {
			if word == w {
				found = true
				break
			}
		}
		if !found {
			diffPenalty += 10
		}
	}

	prefixBonus := 0.0
	minLen := min(len(a), len(b))
	for i := 0; i < minLen && i < len(a) && i < len(b) && a[i] == b[i]; i++ {
		prefixBonus += 2
	}
	if prefixBonus > 10 {
		prefixBonus = 10
	}

	score := basePercent - diffPenalty + prefixBonus
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score / 100
}

func (l *LCSMatcher) AddDocument(text string, answer []int) {
	words := l.tokenize(text)
	fmt.Printf("  📚 [LCS] AddDocument: '%s' -> tokens: %v\n", text, words)
	l.Documents = append(l.Documents, words)
	l.Answers = append(l.Answers, answer)
	fmt.Printf("  📚 [LCS] Now %d documents\n", len(l.Documents))
}

func (l *LCSMatcher) Predict(query string) []int {
	fmt.Printf("\n  🔍 [LCS] Predict called with query: '%s'\n", query)

	if len(l.Documents) == 0 {
		fmt.Printf("  ❌ [LCS] No documents\n")
		return nil
	}

	queryWords := l.tokenize(query)
	fmt.Printf("  🔍 [LCS] Query tokens: %v\n", queryWords)

	bestIdx := -1
	bestScore := -1.0

	for i, doc := range l.Documents {
		score := l.lcsSimilarity(queryWords, doc)
		fmt.Printf("  📊 [LCS] Doc %d: score=%.4f | doc='%v'\n", i, score, doc)
		if score > bestScore {
			bestScore = score
			bestIdx = i
		}
	}

	fmt.Printf("  📊 [LCS] Best score: %.4f at index %d\n", bestScore, bestIdx)
	fmt.Printf("  📊 [LCS] Threshold: %.2f\n", config.LCSConfidenceThreshold)

	if bestIdx >= 0 && bestScore > config.LCSConfidenceThreshold {
		fmt.Printf("  ✅ [LCS] Returning answer at index %d\n", bestIdx)
		return l.Answers[bestIdx]
	}
	fmt.Printf("  ❌ [LCS] Returning nil (score below threshold)\n")
	return nil
}

func (l *LCSMatcher) GetConfidence(query string) float64 {
	if len(l.Documents) == 0 {
		return 0
	}

	queryWords := l.tokenize(query)
	bestScore := -1.0
	for _, doc := range l.Documents {
		score := l.lcsSimilarity(queryWords, doc)
		if score > bestScore {
			bestScore = score
		}
	}
	return bestScore
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
