package config

const (
	ConversationsFile = "data/conversations.txt"
	CheckpointFile    = "data/checkpoint.gob"
)

const (
	SimilarityThreshold = 0.3
	ExactMatchBoost     = 1.5
	MinWordMatchScore   = 0.5
	TopKCandidates      = 3
	TopKTempHigh        = 5
	TopKTempLow         = 1
	IDFSmoothing        = 1.0

	DefaultTemp = 0.7
	MinTemp     = 0.1
	MaxTemp     = 1.5
)
