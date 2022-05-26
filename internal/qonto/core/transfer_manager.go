package core

import (
	"context"
	"fmt"
	"time"

	"github.com/maxim-nazarenko/qonto-interview/internal/qonto/storage"
)

type (
	qontoTransferManager struct {
		storage storage.Storage
	}
)

func NewQontoTransferManager(storage storage.Storage) *qontoTransferManager {
	return &qontoTransferManager{
		storage: storage,
	}
}

func (qm *qontoTransferManager) ProcessTransfers(ctx context.Context, request *Request) error {
	transactions := make([]*storage.Transaction, 0, len(request.CreditTransfers))

	return qm.storage.WithTransactionStorage(ctx, func(ctx context.Context, txStorage storage.Storage) error {
		var totalAmount int64
		for _, ct := range request.CreditTransfers {
			totalAmount += ct.Amount.Cents
		}
		account, err := txStorage.FindAccountByIBAN(ctx, request.Party.IBAN)
		if err != nil {
			return err
		}

		if account.BalanceCents < totalAmount {
			return ErrNotEnoughFunds
		}

		for _, tx := range request.CreditTransfers {
			transactions = append(transactions,
				&storage.Transaction{
					CounterpartyName: tx.CounterParty.Name,
					CounterpartyIBAN: tx.CounterParty.IBAN,
					CounterpartyBIC:  tx.CounterParty.BIC,
					AmountCents:      tx.Amount.Cents,
					AmountCurrency:   string(CURRENCY_EURO),
					BankAccountID:    account.ID,
					Description:      fmt.Sprintf("[%s] Transfer to %s", time.Now().UTC().Format(time.RFC3339), tx.CounterParty.Name),
				},
			)
		}

		if err := txStorage.UpdateAccountBalance(ctx, account.ID, account.BalanceCents-totalAmount); err != nil {
			return err
		}

		if err := txStorage.AppendAccountTransactions(ctx, transactions); err != nil {
			return err
		}

		return nil
	})
}
