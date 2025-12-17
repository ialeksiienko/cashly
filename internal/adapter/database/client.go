package database

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type WithTransaction interface {
	WithTransaction(ctx context.Context, fn func(pgx.Tx) error) error
}
