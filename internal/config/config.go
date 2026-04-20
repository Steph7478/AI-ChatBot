package config

const (
	ModelFile         = "data/model.gob"
	ConversationsFile = "data/conversations.txt"
	TrainingDataFile  = "data/training_data.txt"
	PromptsFile       = "data/prompts.txt"
)

const (
	VocabSize    = 5000
	EmbeddingDim = 64
	HiddenDim    = 128
	NumHeads     = 4
	NumLayers    = 2
	MaxSeqLen    = 20
	Dropout      = 0.2
)

const (
	DefaultTemp   = 1.2
	TopK          = 50
	MaxTokens     = 50
	MinConfidence = 0.2
)

const (
	MinSimilarityScore = 0.8
	ContainBonus       = 0.2
	LengthPenalty      = 0.02
)

const (
	MinTemp = 0.1
	MaxTemp = 2.0
)

const (
	LearningRate   = 0.01
	BatchSize      = 16
	TrainMaxSeqLen = 20
	GradientClip   = 0.5
	Epsilon        = 1e-7
	Patience       = 50
)
