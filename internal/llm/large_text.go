package llm

import (
	"context"
	"log"

	"github.com/pkg/errors"
)

func LargetTextEmbeddings(ctx context.Context, srv Service, text string, maxChunkSize int) ([]float64, error) {
	totalChunks := len(text)/maxChunkSize + 1

	vectors := make([][]float64, 0)

	for i := 0; i < totalChunks; i++ {
		start := i * maxChunkSize
		end := start + maxChunkSize
		if end > len(text) {
			end = len(text)
		}

		chunk := text[start:end]

		embeddings, err := srv.Embeddings(ctx, chunk)
		if err != nil {
			log.Printf("[CHUNK #%d IGNORED] %+v", i, errors.WithStack(err))
			continue
		}

		vectors = append(vectors, embeddings)
	}

	embeddings := average(vectors...)

	return embeddings, nil
}

func average(vectors ...[]float64) []float64 {
	totalVectors := len(vectors)
	vectorLen := len(vectors[0])

	average := make([]float64, vectorLen)

	for i := 0; i < vectorLen; i++ {
		var sum float64
		for _, v := range vectors {
			sum += v[i]
		}
		average[i] = sum / float64(totalVectors)
	}

	return average
}
