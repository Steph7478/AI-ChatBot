package transformer

import (
	"chatbot/internal/neural/layers"
	"encoding/gob"
	"fmt"
	"os"
)

func init() {
	gob.Register(&Transformer{})
	gob.Register(&TransformerBlock{})
	gob.Register(&layers.MultiHeadAttention{})
	gob.Register(&layers.FeedForwardLayer{})
	gob.Register(&layers.EmbeddingLayer{})
	gob.Register(&layers.OutputLayer{})
}

func (t *Transformer) Save(path string) error {
	fmt.Printf("💾 Saving transformer to %s...\n", path)

	if t.Embedding != nil {
		var sum float64
		var count int
		for i := 0; i < 10 && i < len(t.Embedding.Weights); i++ {
			for j := 0; j < 10 && j < len(t.Embedding.Weights[i]); j++ {
				sum += t.Embedding.Weights[i][j]
				count++
			}
		}
	}

	f, err := os.Create(path)
	if err != nil {
		fmt.Printf("❌ Create file error: %v\n", err)
		return err
	}
	defer f.Close()

	if err := gob.NewEncoder(f).Encode(t); err != nil {
		fmt.Printf("❌ Encode error: %v\n", err)
		return err
	}

	fmt.Printf("✅ Transformer saved successfully (%d bytes)\n", getFileSize(path))
	return nil
}

func (t *Transformer) Load(path string) error {
	fmt.Printf("📂 Loading transformer from %s...\n", path)
	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("⚠️  Open file error: %v\n", err)
		return err
	}
	defer f.Close()

	temp := &Transformer{}
	if err := gob.NewDecoder(f).Decode(temp); err != nil {
		fmt.Printf("❌ Decode error: %v\n", err)
		return err
	}

	t.Embedding = temp.Embedding
	t.Blocks = temp.Blocks
	t.Output = temp.Output
	t.MaxSeqLen = temp.MaxSeqLen
	t.Dropout = temp.Dropout

	fmt.Printf("✅ Transformer loaded successfully\n")
	return nil
}

func getFileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}
