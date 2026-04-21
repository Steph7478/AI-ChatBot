package layers

import (
	"chatbot/internal/core"
	"math"
)

func NewEmbeddingLayer(vocabSize, dim int) *EmbeddingLayer {
	return &EmbeddingLayer{
		Weights: core.NewMatrix(vocabSize, dim),
		Dim:     dim,
	}
}

func (e *EmbeddingLayer) Forward(ids []int) [][]float64 {
	if len(ids) == 0 {
		return [][]float64{}
	}
	out := make([][]float64, len(ids))
	for i, id := range ids {
		if id >= 0 && id < len(e.Weights) {
			out[i] = make([]float64, len(e.Weights[id]))
			copy(out[i], e.Weights[id])
		} else {
			out[i] = make([]float64, e.Dim)
		}
	}
	return out
}

var peCache = make(map[int][][]float64)

func PositionalEncoding(seqLen, dim int) [][]float64 {
	if c, ok := peCache[seqLen]; ok && len(c) > 0 && len(c[0]) == dim {
		return c
	}
	pe := core.Zeros(seqLen, dim)
	for i := range seqLen {
		for j := range dim {
			if j%2 == 0 {
				pe[i][j] = math.Sin(float64(i) / math.Pow(10000, float64(j)/float64(dim)))
			} else {
				pe[i][j] = math.Cos(float64(i) / math.Pow(10000, float64(j-1)/float64(dim)))
			}
		}
	}
	peCache[seqLen] = pe
	return pe
}
