package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"oybek.io/kerege/model"
)

type DocumentStorage struct {
	conn *pgx.Conn
}

func NewDocumentStorage(conn *pgx.Conn) *DocumentStorage {
	return &DocumentStorage{
		conn: conn,
	}
}

func (ds *DocumentStorage) Store(ctx context.Context, document *model.Document) error {
	return InsertDocument(ds.conn, ctx, document.Embedding, document.Content)
}

func (ds *DocumentStorage) Search(ctx context.Context, embedding []float32) ([]string, error) {
	return SearchDocuments(ds.conn, ctx, embedding, 1)
}
