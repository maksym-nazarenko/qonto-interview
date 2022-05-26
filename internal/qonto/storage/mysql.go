package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

type (
	WaiterFunc   func(db *sql.DB) (bool, error)
	mysqlStorage struct {
		db      *sql.DB
		querier Querier
	}

	Querier interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
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
		db:      db,
		querier: db,
	}, nil
}

func (m *mysqlStorage) CreateAccount(ctx context.Context, name, iban, bic string, initialBalanceCents int64) (int64, error) {
	stmt := `
		INSERT INTO bank_accounts ( organization_name, balance_cents, iban, bic)
		VALUES (?,?,?,?)`

	result, err := m.querier.ExecContext(ctx, stmt, name, initialBalanceCents, iban, bic)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *mysqlStorage) FindAccount(ctx context.Context, id int64) (Account, error) {
	stmt := `
		SELECT
			id, organization_name, balance_cents, iban, bic
		FROM
			bank_accounts
		WHERE id = ?
		`

	row := m.querier.QueryRowContext(ctx, stmt, id)
	account := Account{}
	if err := row.Scan(&account.ID, &account.Name, &account.BalanceCents, &account.IBAN, &account.BIC); err != nil {
		return Account{}, err
	}
	return account, nil
}

func (m *mysqlStorage) FindAccountByIBAN(ctx context.Context, iban string) (Account, error) {
	stmt := `
		SELECT
			id, organization_name, balance_cents, iban, bic
		FROM
			bank_accounts
		WHERE iban = ?
		`

	row := m.querier.QueryRowContext(ctx, stmt, iban)
	account := Account{}
	if err := row.Scan(&account.ID, &account.Name, &account.BalanceCents, &account.IBAN, &account.BIC); err != nil {
		return Account{}, err
	}
	return account, nil
}

func (m *mysqlStorage) UpdateAccountBalance(ctx context.Context, id, balance int64) error {
	stmt := `
		UPDATE
			bank_accounts
		SET
			balance_cents = ?
		WHERE id = ?
		`

	_, err := m.querier.ExecContext(ctx, stmt, balance, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *mysqlStorage) FindAccountTransactions(ctx context.Context, id int64) ([]*Transaction, error) {
	stmt := `
		SELECT
			id,
			counterparty_name, counterparty_iban, counterparty_bic,
			amount_cents, amount_currency,
			bank_account_id,
			description
		FROM
			transactions
		WHERE bank_account_id = ?
		`

	rows, err := m.querier.QueryContext(ctx, stmt, id)
	if err != nil {
		return nil, err
	}
	result := []*Transaction{}
	defer rows.Close()

	for rows.Next() {
		tx := Transaction{}
		if err := rows.Scan(
			&tx.ID,
			&tx.CounterpartyName, &tx.CounterpartyIBAN, &tx.CounterpartyBIC,
			&tx.AmountCents, &tx.AmountCurrency,
			&tx.BankAccountID,
			&tx.Description,
		); err != nil {
			return nil, err
		}
		result = append(result, &tx)
	}

	return result, rows.Err()
}

func (m *mysqlStorage) AppendAccountTransactions(ctx context.Context, transactions []*Transaction) error {
	stmt := `
		INSERT INTO
			transactions
			(
				counterparty_name,
				counterparty_iban,
				counterparty_bic,
				amount_cents,
				amount_currency,
				bank_account_id,
				description
			)
		VALUES
		` + strings.Repeat(", (?, ?, ?, ?, ?, ?, ?)", len(transactions))[1:]

	args := []interface{}{}
	for _, v := range transactions {
		args = append(args,
			v.CounterpartyName,
			v.CounterpartyIBAN,
			v.CounterpartyBIC,
			v.AmountCents,
			v.AmountCurrency,
			v.BankAccountID,
			v.Description)
	}
	_, err := m.querier.ExecContext(ctx, stmt, args...)
	return err
}

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

func (m *mysqlStorage) WithTransactionStorage(ctx context.Context, f func(context.Context, Storage) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	txMySQL := &mysqlStorage{
		db:      m.db,
		querier: tx,
	}
	if err := f(ctx, txMySQL); err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			return fmt.Errorf("%w: cannot rollback transaction: %v", err, errTx)
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
