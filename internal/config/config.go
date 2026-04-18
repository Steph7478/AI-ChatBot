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
	Temperature   = 1.2
	TopK          = 50
	MaxTokens     = 100
	MinConfidence = 0.1
)

const (
	MinSimilarityScore = 0.8
	ExactWordBonus     = 0.3
	ContainBonus       = 0.2
	LengthPenalty      = 0.03
)

const (
	DefaultTemp = 0.9
	MinTemp     = 0.1
	MaxTemp     = 2.0
)

const (
	LearningRate   = 0.0005
	BatchSize      = 8
	TrainMaxSeqLen = 64
	GradientClip   = 0.5
	Epsilon        = 1e-6
)
