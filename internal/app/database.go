package app

import (
	"cashly/internal/pkg/sl"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DatabaseConfig struct {
	Username string
	Password string
	Hostname string
	Port     string
	DBName   string

	Logger *sl.MyLogger
}

func (db DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		db.Username, db.Password, db.Hostname, db.Port, db.DBName)
}

type Datastore struct {
	dbPool *pgxpool.Pool
}

func NewDatastore(pool *pgxpool.Pool) Datastore {
	return Datastore{dbPool: pool}
}

func (ds Datastore) Pool() *pgxpool.Pool {
	return ds.dbPool
}

func NewDBPool(dbConfig DatabaseConfig) (*pgxpool.Pool, func(), error) {

	f := func() {}

	pool, err := pgxpool.Connect(context.Background(), dbConfig.DSN())
	if err != nil {
		return nil, f, errors.New("database connection error")
	}

	err = validateDBPool(pool, dbConfig.Logger)
	if err != nil {
		return nil, f, err
	}
	return pool, func() { pool.Close() }, nil
}

func validateDBPool(pool *pgxpool.Pool, logger *sl.MyLogger) error {
	err := pool.Ping(context.Background())
	if err != nil {
		return errors.New("database connection error")
	}

	var (
		currentDatabase string
		currentUser     string
		dbVersion       string
	)

	sqlStatement := `select current_database(), current_user, version();`
	row := pool.QueryRow(context.Background(), sqlStatement)
	err = row.Scan(&currentDatabase, &currentUser, &dbVersion)

	switch {
	case err == sql.ErrNoRows:
		return errors.New("no rows were returned")
	case err != nil:
		return errors.New("database connection error")
	default:
		logger.Debug(fmt.Sprintf("database version: %s\n", dbVersion))
		logger.Debug(fmt.Sprintf("current database user: %s\n", currentUser))
		logger.Debug(fmt.Sprintf("current database: %s\n", currentDatabase))
	}

	return nil
}
