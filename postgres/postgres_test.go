package postgres

import (
	"context"
	"os"
	"testing"
	"time"
)

var testdb *TestConn

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := RunPostgres(ctx)
	if err != nil {
		return
	}

	db.Conn.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS vector")
	db.Conn.Exec(ctx, `CREATE TABLE documents (
                         id SERIAL PRIMARY KEY,
                         content TEXT,
                         embedding VECTOR(1536)
                     )`)
	testdb = db

	code := m.Run()
	db.Terminate(ctx)
	os.Exit(code)
}
