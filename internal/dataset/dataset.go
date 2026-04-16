package dataset

import (
	"bufio"
	"os"
	"strings"
	"unicode"
)

type Vocabulary struct {
	WordToID map[string]int
	IDToWord map[int]string
	Size     int
}

func NewVocabulary() *Vocabulary {
	vocab := &Vocabulary{
		WordToID: make(map[string]int),
		IDToWord: make(map[int]string),
	}
	vocab.AddWord("<pad>")
	vocab.AddWord("<unk>")
	vocab.AddWord("<eos>")
	return vocab
}

func (v *Vocabulary) AddWord(word string) int {
	if id, exists := v.WordToID[word]; exists {
		return id
	}
	id := v.Size
	v.WordToID[word] = id
	v.IDToWord[id] = word
	v.Size++
	return id
}

func (v *Vocabulary) GetID(word string) int {
	if id, exists := v.WordToID[word]; exists {
		return id
	}
	return v.WordToID["<unk>"]
}

type Dataset struct {
	Conversations [][2][]int
	Vocab         *Vocabulary
}

func NewDataset(vocab *Vocabulary) *Dataset {
	return &Dataset{
		Conversations: make([][2][]int, 0),
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
		} else if strings.HasPrefix(line, "Bot:") {
			currentBot = strings.TrimSpace(strings.TrimPrefix(line, "Bot:"))
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
	d.Conversations = append(d.Conversations, [2][]int{userTokens, botTokens})
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

func (d *Dataset) Tokenize(text string) []int {
	words := tokenizeText(text)
	tokens := make([]int, 0, len(words))
	for _, word := range words {
		tokens = append(tokens, d.Vocab.GetID(word))
	}
	return tokens
}

func (d *Dataset) Detokenize(tokens []int) string {
	words := make([]string, 0, len(tokens))
	for _, id := range tokens {
		if word, exists := d.Vocab.IDToWord[id]; exists {
			if !strings.HasPrefix(word, "<") {
				words = append(words, word)
			}
		}
	}
	return strings.Join(words, " ")
}

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
