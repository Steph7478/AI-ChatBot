package layers

type MultiHeadAttention struct {
	NumHeads   int
	HeadDim    int
	WQ, WK, WV [][]float64
	WO         [][]float64
}

type EmbeddingLayer struct {
	Weights [][]float64
	Dim     int
}

type FeedForwardLayer struct {
	W1, W2 [][]float64
}

type OutputLayer struct {
	Weights [][]float64
}
