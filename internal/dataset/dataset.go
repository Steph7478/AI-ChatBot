package dataset

import (
	"bufio"
	"os"
	"strings"
	"unicode"
)

func NewDataset(vocab *Vocabulary) *Dataset {
	return &Dataset{
		Conversations: make([]Conversation, 0),
		Vocab:         vocab,
	}
}

func (d *Dataset) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentUser string
	var currentBot string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			if currentUser != "" && currentBot != "" {
				d.AddConversation(currentUser, currentBot)
				currentUser = ""
				currentBot = ""
			}
			continue
		}

		if strings.HasPrefix(line, "User:") {
			if currentUser != "" && currentBot != "" {
				d.AddConversation(currentUser, currentBot)
				currentUser = ""
				currentBot = ""
			}
			currentUser = strings.TrimSpace(strings.TrimPrefix(line, "User:"))
		} else if after, ok := strings.CutPrefix(line, "Bot:"); ok {
			currentBot = strings.TrimSpace(after)
		}
	}

	if currentUser != "" && currentBot != "" {
		d.AddConversation(currentUser, currentBot)
	}

	return scanner.Err()
}

func (d *Dataset) AddConversation(userMsg, botMsg string) {
	userTokens := d.tokenizeAndAdd(userMsg)
	botTokens := d.tokenizeAndAdd(botMsg)
	d.Conversations = append(d.Conversations, Conversation{userTokens, botTokens})
}

func (d *Dataset) tokenizeAndAdd(text string) []int {
	words := tokenizeText(text)
	tokens := make([]int, 0, len(words))
	for _, word := range words {
		id := d.Vocab.AddWord(word)
		tokens = append(tokens, id)
	}
	return tokens
}

// Tokenize converte texto em IDs usando o vocabulário existente
func (d *Dataset) Tokenize(text string) []int {
	words := tokenizeText(text)
	tokens := make([]int, 0, len(words))
	for _, word := range words {
		tokens = append(tokens, d.Vocab.GetID(word))
	}
	return tokens
}

// Detokenize converte IDs de volta para texto
func (d *Dataset) Detokenize(tokens []int) string {
	return d.Vocab.Detokenize(tokens)
}

// tokenizeText função auxiliar para tokenização básica
func tokenizeText(text string) []string {
	text = strings.ToLower(text)

	var result strings.Builder
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' {
			result.WriteRune(r)
		} else {
			result.WriteRune(' ')
		}
	}

	return strings.Fields(result.String())
}

func (d *Dataset) GetConversationCount() int {
	return len(d.Conversations)
}

func (d *Dataset) GetConversation(index int) ([]int, []int, bool) {
	if index < 0 || index >= len(d.Conversations) {
		return nil, nil, false
	}
	return d.Conversations[index][0], d.Conversations[index][1], true
}
