package config

const (
	ModelFile         = "data/model.gob"
	ConversationsFile = "data/conversations.txt"
	TrainingDataFile  = "data/training_data.txt"
	PromptsFile       = "data/prompts.txt"
)

const (
	VocabSize    = 8000
	EmbeddingDim = 96
	HiddenDim    = 192
	NumHeads     = 6
	NumLayers    = 3
	MaxSeqLen    = 25
	Dropout      = 0.2
)

const (
	DefaultTemp   = 0.7
	TopK          = 40
	MaxTokens     = 50
	MinConfidence = 0.2
)

const (
	MinSimilarityScore = 0.75
	ContainBonus       = 0.2
	LengthPenalty      = 0.02
)

const (
	MinTemp = 0.1
	MaxTemp = 2.0
)

const (
	LearningRate   = 0.003
	BatchSize      = 16
	TrainMaxSeqLen = 25
	GradientClip   = 0.5
	Epsilon        = 1e-7
	Patience       = 50
)
