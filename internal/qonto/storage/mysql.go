package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

type (
	WaiterFunc   func(db *sql.DB) (bool, error)
	mysqlStorage struct {
		db *sql.DB
	}

	Querier interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	}
)

// TimeoutPingWaiter waits until the DB is available or fails on timeout
var TimeoutPingWaiter func(context.Context, time.Duration) WaiterFunc = func(parentCtx context.Context, timeout time.Duration) WaiterFunc {
	var err error
	ctx, cancel := context.WithTimeout(parentCtx, timeout)
	return func(db *sql.DB) (bool, error) {
		for {
			select {
			case <-ctx.Done():
				cancel()
				if err != nil {
					err = errors.New(ctx.Err().Error() + ": " + err.Error())
				}
				return false, err
			case <-time.After(1 * time.Second):
				if err = db.PingContext(ctx); err == nil {
					cancel()
					return false, nil
				}
				return true, err
			}
		}
	}
}

// NewMysqlStorage creates and initializes new MySQL storage instance
func NewMysqlStorage(config *mysql.Config) (*mysqlStorage, error) {
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(30)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &mysqlStorage{
		db: db,
	}, nil
}

// WithTransaction implements Storage interface
func (m *mysqlStorage) WithTransaction(ctx context.Context, f func(context.Context, Querier) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := f(ctx, tx); err != nil {
		if errTx := tx.Rollback(); err != nil {
			return fmt.Errorf("cannot rollback transaction: %v, original error: %v", errTx, err)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m *mysqlStorage) Close() error {
	return m.db.Close()
}

func (m *mysqlStorage) DB() *sql.DB {
	return m.db

}

func (m *mysqlStorage) Wait(f WaiterFunc) error {
	cont, err := f(m.db)
	for cont {
		cont, err = f(m.db)
	}

	return err
}

// NewMysqlConfig initializes new MySQL connection configuration with sane defaults
func NewMysqlConfig() *mysql.Config {
	mysqlConfig := mysql.NewConfig()
	mysqlConfig.AllowNativePasswords = true
	mysqlConfig.MultiStatements = true // if false, SQL with >1 statement (e.g. create table in migrations) will fail
	mysqlConfig.ParseTime = true

	return mysqlConfig
}
