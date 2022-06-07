package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/yledovskikh/devops-tpl/internal/storage"
)

const (
	upsertGauge   = "INSERT INTO mgauges (metric_name, metric_value) VALUES($1, $2) ON CONFLICT (metric_name) DO UPDATE SET metric_value = $2 WHERE mgauges.metric_name = $1;"
	upsertCounter = "INSERT INTO mcounter (metric_name, metric_value) VALUES($1, $2) ON CONFLICT (metric_name) DO UPDATE SET metric_value = mcounter.metric_value+$2 WHERE mcounter.metric_name = $1;"
)

type DB struct {
	Pool *pgxpool.Pool
	ctx  context.Context
}

func (d *DB) Close() {
	d.Pool.Close()
}

func New(dsn string, ctx context.Context) (*DB, error) {

	dbPool, err := pgxpool.Connect(context.Background(), dsn)

	if err != nil {
		return &DB{}, err
	}
	err = dbMigrate(dbPool, ctx)
	if err != nil {
		return &DB{}, err
	}
	return &DB{dbPool, ctx}, nil
}

func (d DB) PingDB() error {

	ctx, cancel := context.WithTimeout(d.ctx, 1*time.Second)
	defer cancel()

	if err := d.Pool.Ping(ctx); err != nil {
		return fmt.Errorf("database is down: %w", err)
	}
	return nil
}

func (d *DB) GetGauge(metricName string) (float64, error) {
	var metric float64
	sql := `SELECT metric_value FROM mgauges WHERE metric_name=$1;`
	row := d.Pool.QueryRow(d.ctx, sql, metricName)
	switch err := row.Scan(&metric); err {
	case pgx.ErrNoRows:
		return 0, storage.ErrNotFound
	case nil:
		return metric, nil
	default:
		return 0, err
	}
}

func (d *DB) SetGauge(metricName string, metricValue float64) error {
	log.Debug().Msgf("Set Counter to DB - metricName:%s, metricValue: %f", metricName, metricValue)
	_, err := d.Pool.Exec(d.ctx, upsertGauge, metricName, metricValue)
	return err
}

func (d *DB) GetAllGauges() map[string]float64 {
	metrics := make(map[string]float64)
	sql := "SELECT metric_name, metric_value FROM mgauges"
	rows, err := d.Pool.Query(d.ctx, sql)
	if err != nil {
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var metricName string
		var metricValue float64
		if err = rows.Scan(&metricName, &metricValue); err != nil {
			return nil
		}
		metrics[metricName] = metricValue
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("")
		return nil
	}

	return metrics
}

func (d *DB) GetCounter(metricName string) (int64, error) {
	var metric int64
	sql := `SELECT metric_value FROM mcounter WHERE metric_name=$1;`
	row := d.Pool.QueryRow(d.ctx, sql, metricName)
	switch err := row.Scan(&metric); err {
	case pgx.ErrNoRows:
		return 0, storage.ErrNotFound
	case nil:
		return metric, nil
	default:
		return 0, err

	}
}

func (d *DB) SetCounter(metricName string, metricValue int64) error {
	log.Debug().Msgf("Set Counter to DB - metricName:%s, metricValue: %d", metricName, metricValue)
	_, err := d.Pool.Exec(d.ctx, upsertCounter, metricName, metricValue)
	return err
}
func (d *DB) GetAllCounters() map[string]int64 {
	metrics := make(map[string]int64)
	sql := "SELECT metric_name, metric_value FROM mgauges"
	rows, err := d.Pool.Query(d.ctx, sql)
	if err != nil {
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var metricName string
		var metricValue int64
		if err = rows.Scan(&metricName, &metricValue); err != nil {
			return nil
		}
		metrics[metricName] = metricValue
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("")
		return nil
	}
	return metrics
}

func dbMigrate(d *pgxpool.Pool, ctx context.Context) error {
	execSQL := []string{
		"CREATE SEQUENCE IF NOT EXISTS serial START 1",
		"CREATE TABLE IF NOT EXISTS mcounter(id integer PRIMARY KEY DEFAULT nextval('serial'), metric_name varchar(255) NOT NULL, metric_value bigint NOT NULL, CONSTRAINT mcounter_metric_name_unique UNIQUE (metric_name))",
		"CREATE TABLE IF NOT EXISTS mgauges(id integer PRIMARY KEY DEFAULT nextval('serial'), metric_name varchar(255) NOT NULL, metric_value double precision NOT NULL, CONSTRAINT mgauges_metric_name_unique UNIQUE (metric_name))",
	}

	for _, sql := range execSQL {
		_, err := d.Exec(ctx, sql)
		if err != nil {
			return fmt.Errorf("database schema was not created - %w", err)
		}
	}

	return nil
}

func (d *DB) SetMetrics(metrics *[]storage.Metric) error {
	tx, err := d.Pool.Begin(d.ctx)
	if err != nil {
		return err
	}
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer tx.Rollback(d.ctx)

	for _, metric := range *metrics {
		switch strings.ToLower(metric.MType) {
		case "gauge":
			_, err = tx.Exec(d.ctx, upsertGauge, metric.ID, *metric.Value)
			if err != nil {
				log.Error().Err(err).Msg("")
			}
		case "counter":
			_, err = tx.Exec(d.ctx, upsertCounter, metric.ID, *metric.Delta)
			if err != nil {
				log.Error().Err(err).Msg("")
			}
		}
	}

	err = tx.Commit(d.ctx)
	if err != nil {
		return err
	}
	return nil
}

//func (d *DB) SetMetric(metric *storage.Metric) error {
//	tx, err := d.Pool.Begin(d.ctx)
//	if err != nil {
//		return err
//	}
//	// Rollback is safe to call even if the tx is already closed, so if
//	// the tx commits successfully, this is a no-op
//	defer tx.Rollback(d.ctx)
//
//		switch strings.ToLower(metric.MType) {
//		case "gauge":
//			_, err = tx.Exec(d.ctx, upsertGauge, metric.ID, *metric.Value)
//			if err != nil {
//				log.Error().Err(err).Msg("")
//			}
//		case "counter":
//			_, err = tx.Exec(d.ctx, upsertCounter, metric.ID, *metric.Delta)
//			if err != nil {
//				log.Error().Err(err).Msg("")
//			}
//		}
//
//	err = tx.Commit(d.ctx)
//	if err != nil {
//		return err
//	}
//	return nil
//}
