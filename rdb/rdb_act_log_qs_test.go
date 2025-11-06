package rdb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"oybek.io/sigma/model"
)

func TestActLogQs(t *testing.T) {
	t.Cleanup(func() {
		cleanDatabase(t.Context(), testdb)
	})

	actLog := model.ActLog{
		MessageID: 12,
		UserID:    13,
		Name:      "coding",
		StartTime: time.Now().UTC().Truncate(time.Millisecond),
		EndTime:   time.Time{},
	}

	t.Run("create act log", func(t *testing.T) {
		err := testdb.CreateActLog(t.Context(), actLog)
		assert.NoError(t, err)
	})

	t.Run("find act log", func(t *testing.T) {
		als, err := testdb.FindActLog(t.Context(), FindActLogArg{
			UserID:    13,
			MessageID: 12,
		})
		assert.NoError(t, err)
		assert.Equal(t, []model.ActLog{actLog}, als)
	})

	t.Run("find active act log", func(t *testing.T) {
		als, err := testdb.FindActLog(t.Context(), FindActLogArg{
			UserID: 13,
			Active: true,
		})
		assert.NoError(t, err)
		assert.Equal(t, []model.ActLog{actLog}, als)
	})

	t.Run("update act log", func(t *testing.T) {
		actLog.EndTime = time.Now().UTC().Truncate(time.Millisecond)
		err := testdb.UpdateActLog(t.Context(), actLog)
		assert.NoError(t, err)
	})

	t.Run("find active act log", func(t *testing.T) {
		als, err := testdb.FindActLog(t.Context(), FindActLogArg{
			UserID: 13,
			Active: true,
		})
		assert.NoError(t, err)
		assert.Equal(t, 0, len(als))
	})

	t.Run("find act log", func(t *testing.T) {
		als, err := testdb.FindActLog(t.Context(), FindActLogArg{
			UserID:    13,
			MessageID: 12,
		})
		assert.NoError(t, err)
		assert.Equal(t, []model.ActLog{actLog}, als)
	})
}
