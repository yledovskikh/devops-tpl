package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
)

func PingDB(dsn string) error {

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		err = errors.New("Unable to connect to database:" + err.Error())
		return err
	}
	defer conn.Close(context.Background())
	return nil

}
