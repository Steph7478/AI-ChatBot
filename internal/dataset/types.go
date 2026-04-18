package dataset

type Conversation [2][]int

type Dataset struct {
	Conversations []Conversation
	Vocab         *Vocabulary
}

type Vocabulary struct {
	WordToID map[string]int
	IDToWord map[int]string
	Size     int
}
