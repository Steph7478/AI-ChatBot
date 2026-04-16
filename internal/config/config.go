package config

const (
	// Embedding dimensions
	EmbedDim = 50

	// Response generation
	MaxTokens   = 50
	Temperature = 0.4

	// Similarity thresholds
	SimilarityThreshold = 0.40
	ExactMatchBoost     = 2.0
	TopKCandidates      = 2
	TopKTempHigh        = 3
	TopKTempLow         = 1

	// IDF smoothing
	IDFSmoothing = 1.0

	// Word matching
	MinWordMatchScore = 0.65
)

// File paths
const (
	ConversationsFile = "data/conversations.txt"
	CheckpointFile    = "data/memory_checkpoint.gob"
)

// Temperature limits
const (
	MinTemperature     = 0.1
	MaxTemperature     = 1.5
	DefaultTemperature = 0.4
)

// Response fallbacks
var GreetingWords = []string{"hello", "hi", "hey", "greetings", "sup", "howdy", "good morning", "good afternoon", "good evening"}
var DefaultResponse = "I don't know how to answer that yet please ask me about java"
var GreetingResponse = "hello welcome to the java chatbot how can I help you learn about java today"
