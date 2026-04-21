package transformer

import "chatbot/internal/neural/layers"

type TransformerBlock struct {
	Attention *layers.MultiHeadAttention
	FFN       *layers.FeedForwardLayer
}

type ForwardCache struct {
	EmbeddingOut [][]float64
	BlockOutputs []BlockCache
}

type BlockCache struct {
	AttentionOut [][]float64
	FFNOut       [][]float64
}

type Transformer struct {
	Embedding *layers.EmbeddingLayer
	Blocks    []*TransformerBlock
	Output    *layers.OutputLayer
	MaxSeqLen int
	Dropout   float64
}
