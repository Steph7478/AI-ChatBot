// internal/config/config.go
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
	MinSimilarityScore = 0.4
	ExactWordBonus     = 0.2
	ContainBonus       = 0.1
	LengthPenalty      = 0.05
)

// Commands
const (
	DefaultTemp = 0.7
	MinTemp     = 0.1
	MaxTemp     = 1.5
)

// Training
const (
	LearningRate   = 0.00001
	BatchSize      = 4
	TrainMaxSeqLen = 32
	GradientClip   = 0.1
	Epsilon        = 1e-7
)
