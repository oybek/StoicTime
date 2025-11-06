package rdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"oybek.io/sigma/model"
)

func TestActQs(t *testing.T) {
	t.Cleanup(func() {
		cleanDatabase(t.Context(), testdb)
	})

	act := model.Act{
		ID:     1,
		UserID: 1,
		Name:   "coding",
	}

	t.Run("create act", func(t *testing.T) {
		err := testdb.CreateAct(t.Context(), act)
		assert.NoError(t, err)
	})

	t.Run("create same act should error", func(t *testing.T) {
		err := testdb.CreateAct(t.Context(), act)
		assert.Error(t, err)
	})

	t.Run("find act", func(t *testing.T) {
		acts, err := testdb.FindAct(t.Context(), FindActArg{
			UserID: act.UserID,
			Name:   act.Name,
		})
		assert.NoError(t, err)
		assert.Equal(t, []model.Act{act}, acts)
	})

	t.Run("find act", func(t *testing.T) {
		acts, err := testdb.FindAct(t.Context(), FindActArg{
			UserID: act.UserID,
		})
		assert.NoError(t, err)
		assert.Equal(t, []model.Act{act}, acts)
	})

	t.Run("delete act", func(t *testing.T) {
		err := testdb.DeleteAct(t.Context(), DeleteActArg{
			UserID: act.UserID,
			Name:   act.Name,
		})
		assert.NoError(t, err)
	})

	t.Run("can't find act", func(t *testing.T) {
		acts, err := testdb.FindAct(t.Context(), FindActArg{
			UserID: act.UserID,
			Name:   act.Name,
		})
		assert.NoError(t, err)
		assert.Equal(t, 0, len(acts))
	})
}
