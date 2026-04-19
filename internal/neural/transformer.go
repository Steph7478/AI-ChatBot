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

	output := &OutputLayer{Weights: newMatrix(embedDim, vocabSize)}

	return &Transformer{
		Embedding: NewEmbeddingLayer(vocabSize, embedDim),
		Blocks:    blocks,
		Output:    output,
		MaxSeqLen: maxSeqLen,
	}
}

func (t *Transformer) Forward(ids []int) [][]float64 {
	if len(ids) == 0 {
		return [][]float64{}
	}
	x := t.Embedding.Forward(ids)
	pe := positionalEncoding(len(x), len(x[0]))
	addMatrices(x, pe)

	for _, block := range t.Blocks {
		attnOut := block.Attention.Forward(x)
		addMatrices(attnOut, x)
		x = attnOut
		ffnOut := block.FFN.Forward(x)
		addMatrices(ffnOut, x)
		x = ffnOut
	}

	output := matMul(x, t.Output.Weights)

	t.Output.OutputCache = make([][]float64, len(x))
	for i := range x {
		t.Output.OutputCache[i] = make([]float64, len(x[i]))
		copy(t.Output.OutputCache[i], x[i])
	}
	return output
}

func (t *Transformer) Generate(input string, tokenizer func(string) []int, cfg InferenceConfig) Response {
	ids := tokenizer(input)
	if len(ids) == 0 {
		return Response{Tokens: []Token{}, Confidence: 0}
	}
	logits := t.Forward(ids)
	if len(logits) == 0 || len(logits[len(logits)-1]) == 0 {
		return Response{Tokens: []Token{}, Confidence: 0}
	}
	lastLogits := logits[len(logits)-1]
	for i := range lastLogits {
		lastLogits[i] /= cfg.Temperature
	}
	probs := softmax(lastLogits)
	tokenID := 0
	maxProb := 0.0
	for i, p := range probs {
		if p > maxProb {
			maxProb = p
			tokenID = i
		}
	}
	return Response{
		Tokens:     []Token{{ID: tokenID, Prob: probs[tokenID]}},
		Confidence: probs[tokenID],
	}
}

func (t *Transformer) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	return encoder.Encode(t)
}

func (t *Transformer) Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	return decoder.Decode(t)
}
