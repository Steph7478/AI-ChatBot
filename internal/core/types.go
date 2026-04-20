package core

type TrainingConfig struct {
	LearningRate float64
	BatchSize    int
	Epochs       int
	Patience     int
}

type InferenceConfig struct {
	Temperature float64
	TopK        int
	MaxTokens   int
}

type Token struct {
	ID   int
	Prob float64
}

type Response struct {
	Tokens     []Token
	Confidence float64
}
