package clients

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type Postgres struct {
	conn *pgx.Conn
	ctx  context.Context
}

func NewPostgres(ctx context.Context, configString string) (*Postgres, error) {
	conn, err := pgx.Connect(ctx, configString)

	if err != nil {
		return nil, err
	}

	return &Postgres{
		conn: conn,
		ctx:  ctx,
	}, nil
}

func (db *Postgres) Query(query string, args ...interface{}) (pgx.Rows, error) {
	return db.conn.Query(db.ctx, query, args...)
}

func (db *Postgres) Close() (bool, error) {
	err := db.conn.Close(db.ctx)

	if err != nil {
		return false, err
	}

	return true, nil
}
