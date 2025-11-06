package rdb

import (
	"context"
	"time"

	"oybek.io/sigma/model"
)

type FindActLogArg struct {
	UserID    int64
	MessageID int64
	Active    bool
}

func (r *Rdb) CreateActLog(ctx context.Context, actLog model.ActLog) error {
	_, err := r.c.Exec(
		ctx,
		`INSERT INTO act_log (user_id, message_id, name, start_time, end_time)
		 VALUES ($1, $2, $3, $4, $5)`,
		actLog.UserID, actLog.MessageID, actLog.Name, actLog.StartTime, actLog.EndTime,
	)
	return err
}

func (r *Rdb) FindActLog(ctx context.Context, arg FindActLogArg) ([]model.ActLog, error) {
	query :=
		`SELECT user_id, message_id, name, start_time, end_time
		 FROM act_log
		 WHERE user_id = $1`
	args := []any{arg.UserID}

	if arg.MessageID != 0 {
		query += ` AND message_id = $2`
		args = append(args, arg.MessageID)
	} else if arg.Active == true {
		query += ` AND end_time = $2`
		args = append(args, time.Time{})
	}

	rows, err := r.c.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	als := []model.ActLog{}
	for rows.Next() {
		var al model.ActLog
		if err := rows.Scan(&al.UserID, &al.MessageID, &al.Name, &al.StartTime, &al.EndTime); err != nil {
			return nil, err
		}
		als = append(als, al)
	}

	return als, nil
}

func (r *Rdb) UpdateActLog(ctx context.Context, actLog model.ActLog) error {
	_, err := r.c.Exec(
		ctx,
		`UPDATE act_log SET name = $1, start_time = $2, end_time = $3
		 WHERE user_id = $4 AND message_id = $5`,
		actLog.Name, actLog.StartTime, actLog.EndTime, actLog.UserID, actLog.MessageID,
	)
	return err
}
