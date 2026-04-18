package neural

import (
	"math"
)

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
		return [][]float64{}
	}
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

func (m *MultiHeadAttention) Backward(grad [][]float64, grads map[string][][]float64, layerIdx int) [][]float64 {
	if len(grad) == 0 {
		return grad
	}

	woGrad := matMul(transpose(grad), m.WO)
	grads[attnKey(layerIdx, "wo")] = woGrad

	gradHeads := splitHeads(grad, m.NumHeads, m.HeadDim)

	headGrads := make([][][]float64, m.NumHeads)
	for h := 0; h < m.NumHeads && h < len(gradHeads); h++ {
		headGrads[h] = gradHeads[h]
	}

	combinedGrad := concatHeads(headGrads)

	if len(combinedGrad) == 0 {
		return grad
	}

	wqGrad := matMul(transpose(combinedGrad), m.WQ)
	wkGrad := matMul(transpose(combinedGrad), m.WK)
	wvGrad := matMul(transpose(combinedGrad), m.WV)

	grads[attnKey(layerIdx, "wq")] = wqGrad
	grads[attnKey(layerIdx, "wk")] = wkGrad
	grads[attnKey(layerIdx, "wv")] = wvGrad

	return combinedGrad
}

func splitHeads(mat [][]float64, numHeads, headDim int) [][][]float64 {
	if len(mat) == 0 {
		return [][][]float64{}
	}
	seqLen := len(mat)
	heads := make([][][]float64, numHeads)

	for h := range numHeads {
		heads[h] = make([][]float64, seqLen)
		for i := range seqLen {
			heads[h][i] = make([]float64, headDim)
			for j := range headDim {
				if h*headDim+j < len(mat[i]) {
					heads[h][i][j] = mat[i][h*headDim+j]
				}
			}
		}
	}
	return heads
}

func attentionScore(q, k [][]float64, headDim int) [][]float64 {
	if len(q) == 0 || len(k) == 0 {
		return [][]float64{}
	}
	scores := matMul(q, transpose(k))
	scale := 1.0 / math.Sqrt(float64(headDim))

	for i := range scores {
		for j := range scores[i] {
			scores[i][j] *= scale
		}
	}
	return softmax2d(scores)
}

func concatHeads(heads [][][]float64) [][]float64 {
	if len(heads) == 0 || len(heads[0]) == 0 {
		return [][]float64{}
	}
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

func attnKey(layer int, name string) string {
	return "attn_" + name + "_" + string(rune(layer+48))
}
