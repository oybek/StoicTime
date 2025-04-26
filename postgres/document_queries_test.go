package postgres

import (
	"context"
	_ "embed"
	"strconv"
	"strings"
	tst "testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/embedding.txt
var floatsRaw string

func TestDocumentQueries(t *tst.T) {
	ctx := context.Background()

	t.Run("InsertDocument", func(t *tst.T) {
		embedding, err := parseFloats(floatsRaw)
		assert.NoError(t, err)
		err = InsertDocument(testdb.Conn, ctx, embedding, "привет")
		assert.NoError(t, err)
	})

	t.Run("SearchDocuments", func(t *tst.T) {
		embedding, err := parseFloats(floatsRaw)
		assert.NoError(t, err)
		documents, err := SearchDocuments(testdb.Conn, ctx, embedding, 1)
		assert.NoError(t, err)
		assert.Equal(t, []string{"привет"}, documents)
	})
}

func parseFloats(s string) ([]float32, error) {
	parts := strings.Fields(s)
	var result []float32
	for _, p := range parts {
		f, err := strconv.ParseFloat(p, 32)
		if err != nil {
			return nil, err
		}
		result = append(result, float32(f))
	}
	return result, nil
}
