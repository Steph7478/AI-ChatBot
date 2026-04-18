package neural

type TrainingConfig struct {
	LearningRate float64
	BatchSize    int
	Epochs       int
}

type Trainer struct {
	Model     *Transformer
	Config    TrainingConfig
	Optimizer *Adam
}

type Adam struct {
	LearningRate float64
	Beta1        float64
	Beta2        float64
	Epsilon      float64
	M            map[string][][]float64
	V            map[string][][]float64
	Step         int
}

type FeedForwardLayer struct {
	W1, W2 [][]float64
}

type EmbeddingLayer struct {
	Weights [][]float64
}

type MultiHeadAttention struct {
	NumHeads   int
	HeadDim    int
	WQ, WK, WV [][]float64
	WO         [][]float64
}

type Transformer struct {
	Embedding *EmbeddingLayer
	Blocks    []*TransformerBlock
	Output    *OutputLayer
	MaxSeqLen int
}

type TransformerBlock struct {
	Attention *MultiHeadAttention
	FFN       *FeedForwardLayer
}

type OutputLayer struct {
	Weights     [][]float64
	OutputCache [][]float64
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

func NewAdam(lr float64) *Adam {
	return &Adam{
		LearningRate: lr,
		Beta1:        0.9,
		Beta2:        0.999,
		Epsilon:      1e-8,
		M:            make(map[string][][]float64),
		V:            make(map[string][][]float64),
		Step:         0,
	}
}
