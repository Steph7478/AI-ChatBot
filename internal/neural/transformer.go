package neural

type Transformer struct {
	Embedding *EmbeddingLayer
	Blocks    []*TransformerBlock
	Output    *OutputLayer
	MaxSeqLen int
}

type TransformerBlock struct {
	Attention *MultiHeadAttention
	FFN       *FeedForwardLayer
}

func NewTransformer(vocabSize, embedDim, hiddenDim, numHeads, numLayers, maxSeqLen int, dropout float64) *Transformer {
	blocks := make([]*TransformerBlock, numLayers)
	for i := range numLayers {
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
	return matMul(x, t.Output.Weights)
}

func (t *Transformer) Generate(input string, tokenizer func(string) []int, cfg InferenceConfig) Response {
	ids := tokenizer(input)
	logits := t.Forward(ids)

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
