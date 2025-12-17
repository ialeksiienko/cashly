package database

import (
	"context"

	"github.com/jackc/pgx/v4"
)

// type Client interface {
// 	Begin(context.Context) (pgx.Tx, error)
// 	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
// 	QueryRow(context.Context, string, ...any) pgx.Row
// 	Query(context.Context, string, ...any) (pgx.Rows, error)
// 	Ping(context.Context) error
// 	Close()
// }

type WithTransaction interface {
	WithTransaction(ctx context.Context, fn func(pgx.Tx) error) error
}
