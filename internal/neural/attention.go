package neural

import (
	"math"
)

type MultiHeadAttention struct {
	NumHeads   int
	HeadDim    int
	WQ, WK, WV [][]float64
	WO         [][]float64
}

func NewMultiHeadAttention(embedDim, numHeads int) *MultiHeadAttention {
	headDim := embedDim / numHeads

	return &MultiHeadAttention{
		NumHeads: numHeads,
		HeadDim:  headDim,
		WQ:       newMatrix(embedDim, embedDim),
		WK:       newMatrix(embedDim, embedDim),
		WV:       newMatrix(embedDim, embedDim),
		WO:       newMatrix(embedDim, embedDim),
	}
}

func (m *MultiHeadAttention) Forward(x [][]float64) [][]float64 {
	q := matMul(x, m.WQ)
	k := matMul(x, m.WK)
	v := matMul(x, m.WV)

	qHeads := splitHeads(q, m.NumHeads, m.HeadDim)
	kHeads := splitHeads(k, m.NumHeads, m.HeadDim)
	vHeads := splitHeads(v, m.NumHeads, m.HeadDim)

	headOutputs := make([][][]float64, m.NumHeads)
	for h := 0; h < m.NumHeads; h++ {
		scores := attentionScore(qHeads[h], kHeads[h], m.HeadDim)
		headOutputs[h] = matMul(scores, vHeads[h])
	}

	concat := concatHeads(headOutputs)
	return matMul(concat, m.WO)
}

func splitHeads(mat [][]float64, numHeads, headDim int) [][][]float64 {
	seqLen := len(mat)
	heads := make([][][]float64, numHeads)

	for h := range numHeads {
		heads[h] = make([][]float64, seqLen)
		for i := range seqLen {
			heads[h][i] = make([]float64, headDim)
			for j := range headDim {
				heads[h][i][j] = mat[i][h*headDim+j]
			}
		}
	}
	return heads
}

func attentionScore(q, k [][]float64, headDim int) [][]float64 {
	scores := matMul(q, transpose(k))
	scale := 1.0 / math.Sqrt(float64(headDim))

	for i := range scores {
		for j := range scores[i] {
			scores[i][j] *= scale
		}
	}
	return softmax2d(scores)
}

func transpose(mat [][]float64) [][]float64 {
	if len(mat) == 0 {
		return [][]float64{}
	}
	result := make([][]float64, len(mat[0]))
	for i := range result {
		result[i] = make([]float64, len(mat))
		for j := range mat {
			result[i][j] = mat[j][i]
		}
	}
	return result
}

func concatHeads(heads [][][]float64) [][]float64 {
	seqLen := len(heads[0])
	embedDim := len(heads[0][0]) * len(heads)

	result := make([][]float64, seqLen)
	for i := range seqLen {
		result[i] = make([]float64, embedDim)
		offset := 0
		for h := range heads {
			for j := 0; j < len(heads[h][i]); j++ {
				result[i][offset+j] = heads[h][i][j]
			}
			offset += len(heads[h][i])
		}
	}
	return result
}
