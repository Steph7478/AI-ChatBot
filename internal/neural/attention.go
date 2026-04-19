package neural

import "math"

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
	if len(x) == 0 {
		return x
	}
	q := matMul(x, m.WQ)
	k := matMul(x, m.WK)
	v := matMul(x, m.WV)

	seqLen := len(x)
	embedDim := len(x[0])

	scores := make([][]float64, seqLen)
	for i := 0; i < seqLen; i++ {
		scores[i] = make([]float64, seqLen)
		for j := 0; j < seqLen; j++ {
			sum := 0.0
			for d := 0; d < embedDim; d++ {
				sum += q[i][d] * k[j][d]
			}
			scores[i][j] = sum / math.Sqrt(float64(embedDim))
		}
	}

	for i := 0; i < seqLen; i++ {
		maxVal := scores[i][0]
		for j := 1; j < seqLen; j++ {
			if scores[i][j] > maxVal {
				maxVal = scores[i][j]
			}
		}
		expSum := 0.0
		for j := 0; j < seqLen; j++ {
			scores[i][j] = math.Exp(scores[i][j] - maxVal)
			expSum += scores[i][j]
		}
		for j := 0; j < seqLen; j++ {
			scores[i][j] /= expSum
		}
	}

	output := make([][]float64, seqLen)
	for i := 0; i < seqLen; i++ {
		output[i] = make([]float64, embedDim)
		for j := 0; j < seqLen; j++ {
			for d := 0; d < embedDim; d++ {
				output[i][d] += scores[i][j] * v[j][d]
			}
		}
	}

	return matMul(output, m.WO)
}

func (m *MultiHeadAttention) Backward(grad [][]float64) (gradInput, gradWQ, gradWK, gradWV, gradWO [][]float64) {
	if len(grad) == 0 {
		return grad, nil, nil, nil, nil
	}

	gradWO = matMul(transpose(grad), m.WO)

	gradOutput := matMul(grad, transpose(m.WO))

	return gradOutput, nil, nil, nil, gradWO
}
