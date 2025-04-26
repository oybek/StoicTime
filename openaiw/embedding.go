package openaiw

import (
	"context"
	"log"

	"github.com/sashabaranov/go-openai"
)

func (s *OpenAIWrapper) GetEmbedding(
	ctx context.Context,
	text string,
) ([]float32, error) {
	resp, err := s.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Model: openai.SmallEmbedding3,
		Input: []string{text},
	})
	if err != nil {
		return nil, err
	}

	log.Printf("OpenAI response: %#v", resp)

	embedding := resp.Data[0].Embedding
	return embedding, nil
}
