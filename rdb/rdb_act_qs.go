package rdb

import (
	"context"

	"oybek.io/sigma/model"
)

type FindActArg struct {
	UserID int64
	Name   string
}

type DeleteActArg struct {
	UserID int64
	Name   string
}

func (r *Rdb) CreateAct(ctx context.Context, act model.Act) error {
	_, err := r.c.Exec(ctx, `INSERT INTO act (user_id, name) VALUES ($1, $2)`, act.UserID, act.Name)
	return err
}

func (r *Rdb) FindAct(ctx context.Context, arg FindActArg) ([]model.Act, error) {
	query := `SELECT id, user_id, name FROM act WHERE user_id = $1`
	args := []any{arg.UserID}
	if arg.Name != "" {
		query += ` AND name = $2`
		args = append(args, arg.Name)
	}

	rows, err := r.c.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	acts := []model.Act{}
	for rows.Next() {
		var act model.Act
		if err := rows.Scan(&act.ID, &act.UserID, &act.Name); err != nil {
			return nil, err
		}
		acts = append(acts, act)
	}

	return acts, nil
}

func (r *Rdb) DeleteAct(ctx context.Context, arg DeleteActArg) error {
	_, err := r.c.Exec(ctx, `DELETE FROM act WHERE user_id = $1 AND name = $2`, arg.UserID, arg.Name)
	return err
}
