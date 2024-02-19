package repository

import (
	"context"
	infraerrors "github.com/diegofsousa/rinha-de-backend-q1/internal/infra"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/gommon/log"
)

type PgConnection struct {
	url string
}

func NewPgConnection(url string) *PgConnection {
	return &PgConnection{
		url: url,
	}
}

func (p PgConnection) Connect(ctx context.Context) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, p.url)
	if err != nil {
		log.Fatal("Unable to connect to database: %v\n", err)
		return nil, infraerrors.GenericError
	}

	return conn, nil
}

func (p PgConnection) Close(ctx context.Context, conn *pgx.Conn) error {
	err := conn.Close(ctx)
	if err != nil {
		log.Fatal("Unable to close conn database: %v\n", err)
		return infraerrors.GenericError
	}

	return nil
}
