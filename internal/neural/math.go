package neural

import "math"

func matMul(a, b [][]float64) [][]float64 {
	if len(a) == 0 || len(b) == 0 {
		return [][]float64{}
	}
	if len(a[0]) != len(b) {
		return [][]float64{}
	}
	result := make([][]float64, len(a))
	for i := 0; i < len(result); i++ {
		result[i] = make([]float64, len(b[0]))
		for k := 0; k < len(a[0]); k++ {
			for j := 0; j < len(b[0]); j++ {
				result[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return result
}

func softmax(vec []float64) []float64 {
	if len(vec) == 0 {
		return vec
	}
	maxVal := vec[0]
	for _, v := range vec {
		if v > maxVal {
			maxVal = v
		}
	}
	expSum := 0.0
	result := make([]float64, len(vec))
	for i, v := range vec {
		result[i] = math.Exp(v - maxVal)
		expSum += result[i]
	}
	for i := 0; i < len(result); i++ {
		if expSum > 0 {
			result[i] /= expSum
		}
	}
	return result
}

func softmax2d(mat [][]float64) [][]float64 {
	if len(mat) == 0 {
		return [][]float64{}
	}
	result := make([][]float64, len(mat))
	for i := 0; i < len(mat); i++ {
		result[i] = softmax(mat[i])
	}
	return result
}

func addMatrices(a, b [][]float64) {
	if len(a) == 0 || len(b) == 0 {
		return
	}
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[i]); j++ {
			if i < len(b) && j < len(b[i]) {
				a[i][j] += b[i][j]
			}
		}
	}
}

func transpose(mat [][]float64) [][]float64 {
	if len(mat) == 0 {
		return [][]float64{}
	}
	result := make([][]float64, len(mat[0]))
	for i := 0; i < len(result); i++ {
		result[i] = make([]float64, len(mat))
		for j := 0; j < len(mat); j++ {
			if i < len(mat[j]) {
				result[i][j] = mat[j][i]
			}
		}
	}
	return result
}
