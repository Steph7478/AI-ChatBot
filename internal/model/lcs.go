package model

import (
	"strings"
)

func NewLCSMatcher() *LCSMatcher {
	return &LCSMatcher{
		Documents: make([][]string, 0),
		Answers:   make([][]int, 0),
	}
}

func (l *LCSMatcher) tokenize(text string) []string {
	text = strings.ToLower(text)
	replacer := strings.NewReplacer(
		".", " ", "!", " ", "?", " ",
		",", " ", ";", " ", ":", " ",
		"\"", " ", "'", " ", "(", " ",
		")", " ", "[", " ", "]", " ",
		"{", " ", "}", " ",
	)
	text = replacer.Replace(text)
	return strings.Fields(text)
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
	if len(a) == 0 && len(b) == 0 {
		return 1.0
	}
	lcsLen := l.lcsLength(a, b)
	maxLen := max(len(a), len(b))
	if maxLen == 0 {
		return 0
	}
	return float64(lcsLen) / float64(maxLen)
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

	if bestIdx >= 0 {
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
