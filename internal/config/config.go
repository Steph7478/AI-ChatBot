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
	NumHeads     = 8
	NumLayers    = 3
	MaxSeqLen    = 30
	Dropout      = 0.2
)

const (
	Temperature   = 0.9
	TopK          = 50
	MaxTokens     = 30
	MinConfidence = 0.3
)

const (
	MinSimilarityScore = 0.7
	ContainBonus       = 0.2
	LengthPenalty      = 0.02
)

const (
	DefaultTemp = 0.9
	MinTemp     = 0.1
	MaxTemp     = 2.0
)

const (
	LearningRate   = 0.001
	BatchSize      = 4
	TrainMaxSeqLen = 30
	GradientClip   = 0.5
	Epsilon        = 1e-7
	Patience       = 5
)
