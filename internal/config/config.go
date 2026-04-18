package config

// File paths
const (
	ModelFile         = "data/model.gob"
	ConversationsFile = "data/conversations.txt"
	PromptsFile       = "data/prompts.txt"
	ExamplesFile      = "data/examples.txt"
)

// Model architecture
const (
	VocabSize    = 10000
	EmbeddingDim = 128
	HiddenDim    = 256
	NumHeads     = 4
	NumLayers    = 2
	MaxSeqLen    = 64
	Dropout      = 0.1
)

// Inference
const (
	Temperature   = 0.7
	TopK          = 5
	MaxTokens     = 50
	MinConfidence = 0.5
)

// Matching
const (
	MinSimilarityScore = 0.5
)

// Commands
const (
	DefaultTemp = 0.7
	MinTemp     = 0.1
	MaxTemp     = 1.5
)
