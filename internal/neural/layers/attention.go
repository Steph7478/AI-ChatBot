package layers

import (
	"chatbot/internal/core"
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
		WQ:       core.NewMatrix(embedDim, embedDim),
		WK:       core.NewMatrix(embedDim, embedDim),
		WV:       core.NewMatrix(embedDim, embedDim),
		WO:       core.NewMatrix(embedDim, embedDim),
	}
}

func (m *MultiHeadAttention) Forward(x [][]float64) [][]float64 {
	if len(x) == 0 {
		return x
	}

	q, k, v := core.MatMul(x, m.WQ), core.MatMul(x, m.WK), core.MatMul(x, m.WV)
	n, d := len(x), len(x[0])

	s := core.Zeros(n, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			sum := 0.0
			for t := 0; t < d; t++ {
				sum += q[i][t] * k[j][t]
			}
			s[i][j] = sum / math.Sqrt(float64(d))
		}
	}

	for i := 0; i < n; i++ {
		mx := s[i][0]
		for j := 1; j < n; j++ {
			if s[i][j] > mx {
				mx = s[i][j]
			}
		}
		sum := 0.0
		for j := 0; j < n; j++ {
			s[i][j] = math.Exp(s[i][j] - mx)
			sum += s[i][j]
		}
		for j := 0; j < n; j++ {
			s[i][j] /= sum
		}
	}

	o := core.Zeros(n, d)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			for t := 0; t < d; t++ {
				o[i][t] += s[i][j] * v[j][t]
			}
		}
	}

	return core.MatMul(o, m.WO)
}
