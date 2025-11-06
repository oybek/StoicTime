package rdb

import (
	"context"
	"embed"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
)

type Rdb struct {
	c *pgx.Conn
}

func NewRdb(ctx context.Context, dbURL string) (*Rdb, error) {
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return &Rdb{c: conn}, nil
}

func (r *Rdb) Disconnect(ctx context.Context) {
	r.c.Close(ctx)
}

//go:embed migrations/*.sql
var fs embed.FS

// Migrate - runs migrations against db
func Migrate(dbURL string) {
	driver, err := iofs.New(fs, "migrations")
	if err != nil {
		log.Fatalf("Migration driver initialization error: %s", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", driver, dbURL)
	if err != nil {
		log.Fatalf("Migration initialization error: %s", err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil {
		log.Printf("Migration result: %s\n", err)
	}
}
