package core

import "math"

func Softmax(v []float64) []float64 {
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

func CrossEntropyLoss(logits [][]float64, targets []int) (loss float64, grads [][]float64) {
	seqLen := len(logits)
	if seqLen == 0 || len(targets) == 0 {
		return 0, nil
	}

	grads = Zeros(seqLen, len(logits[0]))
	loss = 0.0

	for t := 0; t < seqLen && t < len(targets); t++ {
		probs := Softmax(logits[t])
		targetIdx := targets[t]

		if targetIdx >= 0 && targetIdx < len(probs) {
			loss -= math.Log(probs[targetIdx] + 1e-10)
			for i := range grads[t] {
				grads[t][i] = probs[i]
			}
			grads[t][targetIdx] -= 1.0
		}
	}
	loss /= float64(min(seqLen, len(targets)))
	return loss, grads
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
