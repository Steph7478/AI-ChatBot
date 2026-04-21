package transformer

import (
	"chatbot/internal/neural/layers"
)

func New(vocabSize, embedDim, hiddenDim, numHeads, numLayers, maxSeqLen int, dropout float64) *Transformer {
	blocks := make([]*TransformerBlock, numLayers)
	for i := range numLayers {
		blocks[i] = &TransformerBlock{
			Attention: layers.NewMultiHeadAttention(embedDim, numHeads),
			FFN:       layers.NewFeedForwardLayer(embedDim, hiddenDim),
		}
	}
	return &Transformer{
		Embedding: layers.NewEmbeddingLayer(vocabSize, embedDim),
		Blocks:    blocks,
		Output:    layers.NewOutputLayer(embedDim, vocabSize),
		MaxSeqLen: maxSeqLen,
		Dropout:   dropout,
	}
}
