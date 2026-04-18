package dataset

import (
	"bufio"
	"os"
	"strings"
)

func NewDataset(vocab *Vocabulary) *Dataset {
	return &Dataset{
		Conversations: make([]Conversation, 0),
		Vocab:         vocab,
	}
}

func (d *Dataset) LoadFromFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var user, bot string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			if user != "" && bot != "" {
				d.AddConversation(user, bot)
				user, bot = "", ""
			}
			continue
		}

		if strings.HasPrefix(line, "User:") {
			if user != "" && bot != "" {
				d.AddConversation(user, bot)
				user, bot = "", ""
			}
			user = strings.TrimSpace(strings.TrimPrefix(line, "User:"))
		} else if after, ok := strings.CutPrefix(line, "Bot:"); ok {
			bot = strings.TrimSpace(after)
		}
	}

	if user != "" && bot != "" {
		d.AddConversation(user, bot)
	}
	return scanner.Err()
}

func (d *Dataset) AddConversation(userMsg, botMsg string) {
	userTokens := d.tokenizeAndAdd(userMsg)
	botTokens := d.tokenizeAndAdd(botMsg)
	d.Conversations = append(d.Conversations, Conversation{userTokens, botTokens})
}

func (d *Dataset) tokenizeAndAdd(text string) []int {
	words := tokenize(text)
	tokens := make([]int, len(words))
	for i, w := range words {
		tokens[i] = d.Vocab.AddWord(w)
	}
	return tokens
}

func (d *Dataset) Tokenize(text string) []int {
	return d.Vocab.Tokenize(text)
}

func (d *Dataset) Detokenize(tokens []int) string {
	return d.Vocab.Detokenize(tokens)
}
