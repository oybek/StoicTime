package service

import (
	"context"

	"oybek.io/kerege/model"
)

type DocumentStorage interface {
	Store(context.Context, *model.Document) error
	Search(context.Context, []float32) ([]string, error)
}

type EmbeddingCalculator interface {
	GetEmbedding(context.Context, string) ([]float32, error)
}
