package service

import (
	"context"

	"oybek.io/sigma/model"
	"oybek.io/sigma/rdb"
)

type ActStorage interface {
	CreateAct(ctx context.Context, act model.Act) error
	FindAct(ctx context.Context, arg rdb.FindActArg) ([]model.Act, error)
	DeleteAct(ctx context.Context, arg rdb.DeleteActArg) error
}

type ActLogStorage interface {
	CreateActLog(ctx context.Context, actLog model.ActLog) error
	FindActLog(ctx context.Context, arg rdb.FindActLogArg) ([]model.ActLog, error)
	UpdateActLog(ctx context.Context, actLog model.ActLog) error
}
