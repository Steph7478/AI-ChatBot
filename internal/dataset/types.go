package dataset

type Conversation [2][]int

type Dataset struct {
	Conversations []Conversation
	Vocab         *Vocabulary
}
