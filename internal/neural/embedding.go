package neural

import (
	"math"
	"math/rand/v2"
)

type EmbeddingLayer struct {
	Weights [][]float64
}

func NewEmbeddingLayer(vocabSize, dim int) *EmbeddingLayer {
	return &EmbeddingLayer{
		Weights: newMatrix(vocabSize, dim),
	}
}

func (e *EmbeddingLayer) Forward(ids []int) [][]float64 {
	embeddings := make([][]float64, len(ids))
	for i, id := range ids {
		if id >= 0 && id < len(e.Weights) {
			embeddings[i] = e.Weights[id]
		} else {
			embeddings[i] = make([]float64, len(e.Weights[0]))
		}
	}
	return embeddings
}

func positionalEncoding(seqLen, dim int) [][]float64 {
	pe := make([][]float64, seqLen)
	for i := 0; i < seqLen; i++ {
		pe[i] = make([]float64, dim)
		for j := 0; j < dim; j++ {
			if j%2 == 0 {
				pe[i][j] = math.Sin(float64(i) / math.Pow(10000, float64(j)/float64(dim)))
			} else {
				pe[i][j] = math.Cos(float64(i) / math.Pow(10000, float64(j-1)/float64(dim)))
			}
		}
	}
	return pe
}

func newMatrix(rows, cols int) [][]float64 {
	m := make([][]float64, rows)
	for i := range m {
		m[i] = make([]float64, cols)
		for j := range m[i] {
			m[i][j] = (rand.Float64() - 0.5) * 0.01
		}
	}
	return m
}
