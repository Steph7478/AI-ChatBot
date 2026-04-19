// internal/config/config.go
package config

const (
	ModelFile         = "data/model.gob"
	ConversationsFile = "data/conversations.txt"
	PromptsFile       = "data/prompts.txt"
)

const (
	VocabSize    = 100
	EmbeddingDim = 16
	HiddenDim    = 32
	NumHeads     = 2
	NumLayers    = 1
	MaxSeqLen    = 32
	Dropout      = 0.1
)

const (
	Temperature   = 0.8
	TopK          = 10
	MaxTokens     = 20
	MinConfidence = 0.3
)

const (
	MinSimilarityScore = 0.35
	ExactWordBonus     = 0.3
	ContainBonus       = 0.2
	LengthPenalty      = 0.02
)

const (
	DefaultTemp = 0.8
	MinTemp     = 0.1
	MaxTemp     = 1.5
)

const (
	LearningRate   = 0.005
	BatchSize      = 8
	TrainMaxSeqLen = 32
	GradientClip   = 0.1
	Epsilon        = 1e-7
)
