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

	tokens := []Token{}
	maxTokens := cfg.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 20
	}

	currentIds := make([]int, len(ids))
	copy(currentIds, ids)

	recentTokens := make(map[int]int)

	for step := 0; step < maxTokens; step++ {
		logits := t.Forward(currentIds)
		if len(logits) == 0 || len(logits[len(logits)-1]) == 0 {
			break
		}

		lastLogits := make([]float64, len(logits[len(logits)-1]))
		copy(lastLogits, logits[len(logits)-1])

		for i := range lastLogits {
			lastLogits[i] /= cfg.Temperature
		}

		probs := softmax(lastLogits)

		// PENALIDADE DE REPETIÇÃO AUMENTADA
		for i := range probs {
			if count, exists := recentTokens[i]; exists && count > 0 {
				// Quanto mais repetiu, menor a chance (penalidade mais forte)
				probs[i] = probs[i] / float64(1+count*3)
			}
		}

		tokenID := 0
		maxProb := 0.0
		for i, p := range probs {
			if p > maxProb {
				maxProb = p
				tokenID = i
			}
		}

		tokens = append(tokens, Token{ID: tokenID, Prob: maxProb})
		recentTokens[tokenID] = recentTokens[tokenID] + 1
		currentIds = append(currentIds, tokenID)

		if len(currentIds) > t.MaxSeqLen {
			currentIds = currentIds[1:]
		}

		if tokenID == 0 {
			break
		}
	}

	if len(tokens) == 0 {
		return Response{Tokens: []Token{}, Confidence: 0}
	}

	confidence := tokens[len(tokens)-1].Prob
	for _, tkn := range tokens {
		if tkn.Prob < confidence {
			confidence = tkn.Prob
		}
	}

	return Response{
		Tokens:     tokens,
		Confidence: confidence,
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
