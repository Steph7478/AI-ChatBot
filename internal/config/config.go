package config

const (
	ModelFile         = "data/model.gob"
	ConversationsFile = "data/conversations.txt"
	PromptsFile       = "data/prompts.txt"
)

const (
	VocabSize    = 10000
	EmbeddingDim = 128
	HiddenDim    = 256
	NumHeads     = 4
	NumLayers    = 2
	MaxSeqLen    = 64
	Dropout      = 0.1
)

const (
	Temperature   = 0.7
	TopK          = 5
	MaxTokens     = 50
	MinConfidence = 0.5
)

const (
	MinSimilarityScore = 0.8
	ExactWordBonus     = 0.2
	ContainBonus       = 0.1
	LengthPenalty      = 0.05
)

const (
	DefaultTemp = 0.7
	MinTemp     = 0.1
	MaxTemp     = 1.5
)

const (
	LearningRate   = 0.00001
	BatchSize      = 4
	TrainMaxSeqLen = 32
	GradientClip   = 0.1
	Epsilon        = 1e-7
)
