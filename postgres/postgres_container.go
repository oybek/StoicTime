package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestConn struct {
	Conn      *pgx.Conn
	Container *testcontainers.Container
}

func RunPostgres(ctx context.Context) (*TestConn, error) {
	user, pass, database := "user", "pass", "db"

	req := testcontainers.ContainerRequest{
		Image:        "ankane/pgvector",
		ExposedPorts: []string{"5432/tcp"},
		AutoRemove:   true,
		Env: map[string]string{
			"POSTGRES_USER":     user,
			"POSTGRES_PASSWORD": pass,
			"POSTGRES_DB":       database,
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	host, err := postgres.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := postgres.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, pass, host, port.Port(), database)
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return &TestConn{Conn: conn, Container: &postgres}, err
}

func (testDB *TestConn) Terminate(ctx context.Context) {
	testDB.Conn.Close(ctx)
	(*testDB.Container).Terminate(ctx)
}
