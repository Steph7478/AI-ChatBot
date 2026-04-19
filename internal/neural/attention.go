package neural

import "math"

func NewMultiHeadAttention(embedDim, numHeads int) *MultiHeadAttention {
	return &MultiHeadAttention{
		NumHeads: numHeads,
		HeadDim:  embedDim / numHeads,
		WQ:       newMatrix(embedDim, embedDim),
		WK:       newMatrix(embedDim, embedDim),
		WV:       newMatrix(embedDim, embedDim),
		WO:       newMatrix(embedDim, embedDim),
	}
}

func (m *MultiHeadAttention) Forward(x [][]float64) [][]float64 {
	if len(x) == 0 {
		return x
	}
	q, k, v := matMul(x, m.WQ), matMul(x, m.WK), matMul(x, m.WV)
	n, d := len(x), len(x[0])
	s := zeros(n, n)
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
	o := zeros(n, d)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			for t := 0; t < d; t++ {
				o[i][t] += s[i][j] * v[j][t]
			}
		}
	}
	return matMul(o, m.WO)
}
