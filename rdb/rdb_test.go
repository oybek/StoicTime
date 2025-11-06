package rdb

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

var testdb *Rdb

// test main runs before any test on current package
func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, pass, db := "user", "pass", "db"
	pgc, err := RunPostgres(ctx, user, pass, db)
	if err != nil {
		panic(err)
	}
	defer pgc.Terminate(ctx)

	host, _ := pgc.Host(ctx)
	port, _ := pgc.MappedPort(ctx, "5432")
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port.Port(), db)

	Migrate(dbURL)
	testdb, _ = NewRdb(ctx, dbURL)

	code := m.Run()
	os.Exit(code)
}

// cleanDatabase truncates all tables in the public schema
func cleanDatabase(ctx context.Context, r *Rdb) error {
	query := `
		DO $$ DECLARE
		    r RECORD;
		BEGIN
		    FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
		        EXECUTE 'TRUNCATE TABLE public.' || quote_ident(r.tablename) || ' RESTART IDENTITY CASCADE';
		    END LOOP;
		END $$;
	`
	_, err := r.c.Exec(ctx, query)
	return err
}
