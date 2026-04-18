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
	for i := range result {
		result[i] = make([]float64, len(b[0]))
		for k := 0; k < len(a[0]); k++ {
			if k >= len(b) {
				continue
			}
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

	for i := range result {
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
	for i := range mat {
		result[i] = softmax(mat[i])
	}
	return result
}

func addMatrices(a, b [][]float64) {
	if len(a) == 0 || len(b) == 0 {
		return
	}
	for i := range a {
		if i >= len(b) {
			continue
		}
		for j := range a[i] {
			if j >= len(b[i]) {
				continue
			}
			a[i][j] += b[i][j]
		}
	}
}

func transpose(mat [][]float64) [][]float64 {
	if len(mat) == 0 {
		return [][]float64{}
	}
	result := make([][]float64, len(mat[0]))
	for i := range result {
		result[i] = make([]float64, len(mat))
		for j := range mat {
			if i < len(mat[j]) {
				result[i][j] = mat[j][i]
			}
		}
	}
	return result
}
