package neural

import "math"

func matMul(a, b [][]float64) [][]float64 {
	if len(a) == 0 || len(b) == 0 || len(a[0]) != len(b) {
		return [][]float64{}
	}
	r := make([][]float64, len(a))
	for i := range a {
		r[i] = make([]float64, len(b[0]))
		for k := range b {
			if a[i][k] == 0 {
				continue
			}
			for j := range b[0] {
				r[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return r
}

func softmax(v []float64) []float64 {
	if len(v) == 0 {
		return v
	}
	m := v[0]
	for _, x := range v[1:] {
		if x > m {
			m = x
		}
	}
	s := 0.0
	r := make([]float64, len(v))
	for i, x := range v {
		r[i] = math.Exp(x - m)
		s += r[i]
	}
	for i := range r {
		r[i] /= s
	}
	return r
}

func addMatrices(a, b [][]float64) {
	for i := range a {
		for j := range a[i] {
			if j < len(b[i]) {
				a[i][j] += b[i][j]
			}
		}
	}
}

func transpose(m [][]float64) [][]float64 {
	if len(m) == 0 {
		return [][]float64{}
	}
	r := make([][]float64, len(m[0]))
	for i := range r {
		r[i] = make([]float64, len(m))
		for j := range m {
			r[i][j] = m[j][i]
		}
	}
	return r
}

func zeros(rows, cols int) [][]float64 {
	r := make([][]float64, rows)
	for i := range r {
		r[i] = make([]float64, cols)
	}
	return r
}

func copyMatrix(src [][]float64) [][]float64 {
	dst := zeros(len(src), len(src[0]))
	for i := range src {
		copy(dst[i], src[i])
	}
	return dst
}

func scale(mat [][]float64, s float64) {
	for i := range mat {
		for j := range mat[i] {
			mat[i][j] *= s
		}
	}
}

func add(a, b [][]float64) [][]float64 {
	r := zeros(len(a), len(a[0]))
	for i := range a {
		for j := range a[i] {
			r[i][j] = a[i][j] + b[i][j]
		}
	}
	return r
}
