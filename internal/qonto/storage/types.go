package storage

import "context"

type (
	Account struct {
		ID           int64
		Name         string
		BalanceCents int64
		BIC          string
		IBAN         string
	}

	Transaction struct {
		ID               int64
		CounterpartyName string
		CounterpartyIBAN string
		CounterpartyBIC  string
		AmountCents      int64
		AmountCurrency   string
		BankAccountID    int64
		Description      string
	}

	// Storage defines interface to be satisfied by concrete storage implementation
	Storage interface {
		// WithTransaction wraps functions in transaction and rolls it back if function returns error
		WithTransaction(context.Context, func(context.Context, Querier) error) error
		WithTransactionStorage(context.Context, func(context.Context, Storage) error) error

		CreateAccount(ctx context.Context, name, iban, bic string, initialBalanceCents int64) (int64, error)
		FindAccount(ctx context.Context, id int64) (Account, error)
		UpdateAccountBalance(ctx context.Context, id, balance int64) error
		FindAccountByIBAN(ctx context.Context, iban string) (Account, error)

		FindAccountTransactions(ctx context.Context, id int64) ([]*Transaction, error)
		AppendAccountTransactions(ctx context.Context, transactions []*Transaction) error

		// Wait runs provided wait function until it returns true without error
		Wait(f WaiterFunc) error

		// Close closes underlying storage connection if supported by concrete implementation
		// Close() error
	}
)
