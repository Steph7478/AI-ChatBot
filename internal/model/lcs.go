package model

import (
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
	l.Documents = append(l.Documents, words)
	l.Answers = append(l.Answers, answer)
}

func (l *LCSMatcher) Predict(query string) []int {
	if len(l.Documents) == 0 {
		return nil
	}

	queryWords := l.tokenize(query)
	bestIdx := -1
	bestScore := -1.0

	for i, doc := range l.Documents {
		score := l.lcsSimilarity(queryWords, doc)
		if score > bestScore {
			bestScore = score
			bestIdx = i
		}
	}

	if bestIdx >= 0 && bestScore > config.LCSConfidenceThreshold {
		return l.Answers[bestIdx]
	}
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
