package layers

import (
	"chatbot/internal/core"
	"math"
)

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
	for i := range n {
		for j := range n {
			sum := 0.0
			for t := range d {
				sum += q[i][t] * k[j][t]
			}
			s[i][j] = sum / math.Sqrt(float64(d))
		}
	}

	for i := range n {
		mx := s[i][0]
		for j := 1; j < n; j++ {
			if s[i][j] > mx {
				mx = s[i][j]
			}
		}
		sum := 0.0
		for j := range n {
			s[i][j] = math.Exp(s[i][j] - mx)
			sum += s[i][j]
		}
		for j := 0; j < n; j++ {
			s[i][j] /= sum
		}
	}

	o := core.Zeros(n, d)
	for i := range n {
		for j := range n {
			for t := range d {
				o[i][t] += s[i][j] * v[j][t]
			}
		}
	}

	return core.MatMul(o, m.WO)
}
