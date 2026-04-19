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
		for j := range n {
			s[i][j] /= sum
		}
	}
	o := zeros(n, d)
	for i := range n {
		for j := range n {
			for t := range d {
				o[i][t] += s[i][j] * v[j][t]
			}
		}
	}
	return matMul(o, m.WO)
}
