package db

import (
	"context"
	//"database/sql"
	"errors"
	"time"

	//"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	//_ "github.com/jackc/pgx/v4/stdlib"
)

type DB struct {
	Pool *pgxpool.Pool
}

func (d *DB) Close() {
	d.Pool.Close()
}

func New(dsn string) (*DB, error) {

	dbPool, err := pgxpool.Connect(context.Background(), dsn)

	if err != nil {
		return &DB{}, err
	}
	return &DB{dbPool}, nil
}

func (d DB) PingDB(ctx context.Context) error {

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := d.Pool.Ping(ctx); err != nil {
		err = errors.New("Database to down:" + err.Error())
		return err
	}
	return nil
}
