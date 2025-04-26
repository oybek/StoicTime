package postgres

import (
	"context"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

func InsertDocument(
	conn *pgx.Conn,
	ctx context.Context,
	embedding []float32,
	content string,
) error {
	embedStr := vectorToPG(embedding)
	_, err := conn.Exec(
		ctx,
		"INSERT INTO documents (content, embedding) VALUES ($1, $2)",
		content, embedStr,
	)
	return err
}

func SearchDocuments(
	conn *pgx.Conn,
	ctx context.Context,
	embedding []float32,
	topK int,
) ([]string, error) {
	embedStr := vectorToPG(embedding)

	rows, err := conn.Query(
		ctx,
		"SELECT content FROM documents ORDER BY embedding <#> $1 LIMIT $2",
		embedStr, topK,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			return nil, err
		}
		results = append(results, content)
	}
	return results, nil
}

func vectorToPG(v []float32) string {
	parts := make([]string, len(v))
	for i, val := range v {
		parts[i] = strconv.FormatFloat(float64(val), 'f', -1, 32)
	}
	return "[" + strings.Join(parts, ",") + "]"
}
