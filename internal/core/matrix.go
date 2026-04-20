package core

import "math/rand/v2"

func MatMul(a, b [][]float64) [][]float64 {
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

func Transpose(m [][]float64) [][]float64 {
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

func Zeros(rows, cols int) [][]float64 {
	if rows == 0 || cols == 0 {
		return [][]float64{}
	}
	r := make([][]float64, rows)
	for i := range r {
		r[i] = make([]float64, cols)
	}
	return r
}

func Copy(src [][]float64) [][]float64 {
	if len(src) == 0 {
		return [][]float64{}
	}
	dst := Zeros(len(src), len(src[0]))
	for i := range src {
		copy(dst[i], src[i])
	}
	return dst
}

func Add(a, b [][]float64) {
	for i := range a {
		for j := range a[i] {
			if j < len(b[i]) {
				a[i][j] += b[i][j]
			}
		}
	}
}

func NewMatrix(rows, cols int) [][]float64 {
	if rows == 0 || cols == 0 {
		return [][]float64{}
	}
	m := Zeros(rows, cols)
	for i := range m {
		for j := range m[i] {
			m[i][j] = (rand.Float64() - 0.5) * 0.01
		}
	}
	return m
}

func AddMatrices(a, b [][]float64) [][]float64 {
	if len(a) == 0 || len(b) == 0 {
		return [][]float64{}
	}
	r := Zeros(len(a), len(a[0]))
	for i := range a {
		for j := range a[i] {
			if j < len(b[i]) {
				r[i][j] = a[i][j] + b[i][j]
			}
		}
	}
	return r
}

func Scale(mat [][]float64, s float64) {
	for i := range mat {
		for j := range mat[i] {
			mat[i][j] *= s
		}
	}
}
