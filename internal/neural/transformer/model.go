package transformer

import (
	"chatbot/internal/neural/layers"
	"fmt"
)

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

func New(vocabSize, embedDim, hiddenDim, numHeads, numLayers, maxSeqLen int, dropout float64) *Transformer {
	blocks := make([]*TransformerBlock, numLayers)
	for i := 0; i < numLayers; i++ {
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

func (t *Transformer) PrintWeightStats() {
	fmt.Println("\n=== MODEL WEIGHT STATS ===")
	if t.Embedding != nil && len(t.Embedding.Weights) > 0 && len(t.Embedding.Weights[0]) > 0 {
		fmt.Printf("Embedding[0][0]: %.6f\n", t.Embedding.Weights[0][0])
		fmt.Printf("Embedding[0][1]: %.6f\n", t.Embedding.Weights[0][1])

		var sum float64
		count := 0
		for i := 0; i < 10 && i < len(t.Embedding.Weights); i++ {
			for j := 0; j < 10 && j < len(t.Embedding.Weights[i]); j++ {
				sum += t.Embedding.Weights[i][j]
				count++
			}
		}
		if count > 0 {
			fmt.Printf("Embedding avg (10x10): %.6f\n", sum/float64(count))
		}
	}

	if t.Output != nil && len(t.Output.Weights) > 0 && len(t.Output.Weights[0]) > 0 {
		fmt.Printf("Output[0][0]: %.6f\n", t.Output.Weights[0][0])
	}
	fmt.Println("==========================")
}
