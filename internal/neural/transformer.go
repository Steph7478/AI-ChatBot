package neural

import (
	"encoding/gob"
	"os"
)

func init() {
	gob.Register(&Transformer{})
	gob.Register(&TransformerBlock{})
	gob.Register(&MultiHeadAttention{})
	gob.Register(&FeedForwardLayer{})
	gob.Register(&EmbeddingLayer{})
	gob.Register(&OutputLayer{})
}

func NewTransformer(vocabSize, embedDim, hiddenDim, numHeads, numLayers, maxSeqLen int, dropout float64) *Transformer {
	blocks := make([]*TransformerBlock, numLayers)
	for i := 0; i < numLayers; i++ {
		blocks[i] = &TransformerBlock{
			Attention: NewMultiHeadAttention(embedDim, numHeads),
			FFN:       NewFeedForwardLayer(embedDim, hiddenDim),
		}
	}
	return &Transformer{
		Embedding: NewEmbeddingLayer(vocabSize, embedDim),
		Blocks:    blocks,
		Output:    &OutputLayer{Weights: newMatrix(embedDim, vocabSize)},
		MaxSeqLen: maxSeqLen,
	}
}

func (t *Transformer) Forward(ids []int) [][]float64 {
	if len(ids) == 0 {
		return [][]float64{}
	}
	x := t.Embedding.Forward(ids)
	addMatrices(x, positionalEncoding(len(x), len(x[0])))
	for _, b := range t.Blocks {
		a := b.Attention.Forward(x)
		addMatrices(a, x)
		x = a
		f := b.FFN.Forward(x)
		addMatrices(f, x)
		x = f
	}
	return matMul(x, t.Output.Weights)
}

func (t *Transformer) Generate(input string, tokenizer func(string) []int, cfg InferenceConfig) Response {
	ids := tokenizer(input)
	if len(ids) == 0 {
		return Response{}
	}
	toks := []Token{}
	cur := make([]int, len(ids))
	copy(cur, ids)
	recent := make(map[int]int)
	mt := cfg.MaxTokens
	if mt <= 0 {
		mt = 20
	}
	for step := 0; step < mt; step++ {
		logits := t.Forward(cur)
		if len(logits) == 0 || len(logits[len(logits)-1]) == 0 {
			break
		}
		last := make([]float64, len(logits[len(logits)-1]))
		copy(last, logits[len(logits)-1])
		for i := range last {
			last[i] /= cfg.Temperature
		}
		p := softmax(last)
		for i := range p {
			if c, ok := recent[i]; ok && c > 0 {
				p[i] = p[i] / float64(1+c*3)
			}
		}
		id, mp := 0, 0.0
		for i, v := range p {
			if v > mp {
				mp, id = v, i
			}
		}
		toks = append(toks, Token{ID: id, Prob: mp})
		recent[id] = recent[id] + 1
		cur = append(cur, id)
		if len(cur) > t.MaxSeqLen {
			cur = cur[1:]
		}
		if id == 0 {
			break
		}
	}
	if len(toks) == 0 {
		return Response{}
	}
	c := toks[len(toks)-1].Prob
	for _, t := range toks {
		if t.Prob < c {
			c = t.Prob
		}
	}
	return Response{Tokens: toks, Confidence: c}
}

func (t *Transformer) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewEncoder(f).Encode(t)
}

func (t *Transformer) Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewDecoder(f).Decode(t)
}
