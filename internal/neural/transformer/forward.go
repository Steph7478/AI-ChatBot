package transformer

import (
	"chatbot/internal/core"
	"chatbot/internal/neural/layers"
)

func (t *Transformer) Forward(ids []int) [][]float64 {
	if len(ids) == 0 {
		return [][]float64{}
	}

	x := t.Embedding.Forward(ids)
	pe := layers.PositionalEncoding(len(x), len(x[0]))
	core.Add(x, pe)

	for _, b := range t.Blocks {
		attnOut := b.Attention.Forward(x)
		core.Add(attnOut, x)
		x = attnOut

		ffnOut := b.FFN.Forward(x)
		core.Add(ffnOut, x)
		x = ffnOut
	}

	return t.Output.Forward(x)
}

func (t *Transformer) ForwardWithCache(ids []int) (logits [][]float64, cache *ForwardCache) {
	if len(ids) == 0 {
		return nil, nil
	}

	cache = &ForwardCache{
		EmbeddingOut: t.Embedding.Forward(ids),
		BlockOutputs: make([]BlockCache, len(t.Blocks)),
	}

	pe := layers.PositionalEncoding(len(cache.EmbeddingOut), t.Embedding.Dim)
	core.Add(cache.EmbeddingOut, pe)
	x := cache.EmbeddingOut

	for i, b := range t.Blocks {
		attnOut := b.Attention.Forward(x)
		core.Add(attnOut, x)
		cache.BlockOutputs[i].AttentionOut = core.Copy(attnOut)

		ffnOut := b.FFN.Forward(attnOut)
		core.Add(ffnOut, attnOut)
		cache.BlockOutputs[i].FFNOut = core.Copy(ffnOut)
		x = ffnOut
	}

	logits = t.Output.Forward(x)
	return logits, cache
}
