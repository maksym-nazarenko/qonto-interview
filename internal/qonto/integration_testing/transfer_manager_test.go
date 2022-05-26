package core

import (
	"context"
	"testing"
	"time"

	"github.com/maxim-nazarenko/qonto-interview/internal/qonto/core"
	"github.com/maxim-nazarenko/qonto-interview/internal/qonto/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessTransfers_happy(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	mysqlStorage, dbName := storage.NewTestDatabase(ctx, t)
	defer mysqlStorage.Close()
	t.Logf("test db name: %s", dbName)

	qontoAccount := core.Party{
		Name: "Qonto customer corp",
		BIC:  "ARWKDJFU",
		IBAN: "UA9935420810036209081725212",
	}

	var accountBalance int64 = 20000
	qontoAccountID, err := mysqlStorage.CreateAccount(ctx, qontoAccount.Name, qontoAccount.IBAN, qontoAccount.BIC, accountBalance)
	require.NoError(t, err)

	transferManager := core.NewQontoTransferManager(mysqlStorage)
	// total amount should be less than accountBalance
	request := core.Request{
		Party: qontoAccount,
		CreditTransfers: []core.Transfer{
			{
				Amount:      core.Amount{Cents: 8000},
				Currency:    core.CURRENCY_EURO,
				Description: "",
				CounterParty: core.Party{
					Name: "counterparty 1",
					BIC:  "bic1",
					IBAN: "iban1",
				},
			},
			{
				Amount:      core.Amount{Cents: 5000},
				Currency:    core.CURRENCY_EURO,
				Description: "",
				CounterParty: core.Party{
					Name: "counterparty 2",
					BIC:  "bic2",
					IBAN: "iban2",
				},
			},
			{
				Amount:      core.Amount{Cents: 3000},
				Currency:    core.CURRENCY_EURO,
				Description: "",
				CounterParty: core.Party{
					Name: "counterparty 3",
					BIC:  "bic3",
					IBAN: "iban3",
				},
			},
			{
				Amount:      core.Amount{Cents: 1000},
				Currency:    core.CURRENCY_EURO,
				Description: "",
				CounterParty: core.Party{
					Name: "counterparty 4",
					BIC:  "bic4",
					IBAN: "iban4",
				},
			},
		},
	}

	err = transferManager.ProcessTransfers(ctx, &request)
	require.NoError(t, err)
	qontoAccountAfterProcessing, err := mysqlStorage.FindAccount(ctx, qontoAccountID)
	require.NoError(t, err)
	var expectedBalance int64 = 3000
	assert.Equal(t, expectedBalance, qontoAccountAfterProcessing.BalanceCents)

	transactions, err := mysqlStorage.FindAccountTransactions(ctx, qontoAccountID)
	require.NoError(t, err)
	assert.Equal(t, len(request.CreditTransfers), len(transactions))
	var expectedTransactionHistoryAmount int64 = 17000
	var actualTransactionsAmount int64
	for _, tx := range transactions {
		actualTransactionsAmount += tx.AmountCents
	}

	assert.Equal(t, expectedTransactionHistoryAmount, actualTransactionsAmount)
}

func TestProcessTransfers_happyAllMoney(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	mysqlStorage, dbName := storage.NewTestDatabase(ctx, t)
	defer mysqlStorage.Close()
	t.Logf("test db name: %s", dbName)

	qontoAccount := core.Party{
		Name: "Qonto customer corp",
		BIC:  "ARWKDJFU",
		IBAN: "UA9935420810036209081725212",
	}

	var accountBalance int64 = 10000
	qontoAccountID, err := mysqlStorage.CreateAccount(ctx, qontoAccount.Name, qontoAccount.IBAN, qontoAccount.BIC, accountBalance)
	require.NoError(t, err)

	transferManager := core.NewQontoTransferManager(mysqlStorage)
	// total amount should be less than accountBalance
	request := core.Request{
		Party: qontoAccount,
		CreditTransfers: []core.Transfer{
			{
				Amount:      core.Amount{Cents: 2000},
				Currency:    core.CURRENCY_EURO,
				Description: "",
				CounterParty: core.Party{
					Name: "counterparty 1",
					BIC:  "bic1",
					IBAN: "iban1",
				},
			},
			{
				Amount:      core.Amount{Cents: 5000},
				Currency:    core.CURRENCY_EURO,
				Description: "",
				CounterParty: core.Party{
					Name: "counterparty 2",
					BIC:  "bic2",
					IBAN: "iban2",
				},
			},
			{
				Amount:      core.Amount{Cents: 3000},
				Currency:    core.CURRENCY_EURO,
				Description: "",
				CounterParty: core.Party{
					Name: "counterparty 3",
					BIC:  "bic3",
					IBAN: "iban3",
				},
			},
		},
	}

	err = transferManager.ProcessTransfers(ctx, &request)
	require.NoError(t, err)
	qontoAccountAfterProcessing, err := mysqlStorage.FindAccount(ctx, qontoAccountID)
	require.NoError(t, err)
	var expectedBalance int64 = 0
	assert.Equal(t, expectedBalance, qontoAccountAfterProcessing.BalanceCents)

	transactions, err := mysqlStorage.FindAccountTransactions(ctx, qontoAccountID)
	require.NoError(t, err)
	assert.Equal(t, len(request.CreditTransfers), len(transactions))
	var expectedTransactionHistoryAmount int64 = accountBalance
	var actualTransactionsAmount int64
	for _, tx := range transactions {
		actualTransactionsAmount += tx.AmountCents
	}

	assert.Equal(t, expectedTransactionHistoryAmount, actualTransactionsAmount)
}

func TestProcessTransfers_decline(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	mysqlStorage, dbName := storage.NewTestDatabase(ctx, t)
	defer mysqlStorage.Close()
	t.Logf("test db name: %s", dbName)

	qontoAccount := core.Party{
		Name: "Qonto customer corp",
		BIC:  "ARWKDJFU",
		IBAN: "UA9935420810036209081725212",
	}

	var accountBalance int64 = 20000
	qontoAccountID, err := mysqlStorage.CreateAccount(ctx, qontoAccount.Name, qontoAccount.IBAN, qontoAccount.BIC, accountBalance)
	require.NoError(t, err)

	transferManager := core.NewQontoTransferManager(mysqlStorage)
	// total amount should be less than accountBalance
	request := core.Request{
		Party: qontoAccount,
		CreditTransfers: []core.Transfer{
			{
				Amount:      core.Amount{Cents: 9000},
				Currency:    core.CURRENCY_EURO,
				Description: "",
				CounterParty: core.Party{
					Name: "counterparty 1",
					BIC:  "bic1",
					IBAN: "iban1",
				},
			},
			{
				Amount:      core.Amount{Cents: 8000},
				Currency:    core.CURRENCY_EURO,
				Description: "",
				CounterParty: core.Party{
					Name: "counterparty 2",
					BIC:  "bic2",
					IBAN: "iban2",
				},
			},
			{
				Amount:      core.Amount{Cents: 3001},
				Currency:    core.CURRENCY_EURO,
				Description: "",
				CounterParty: core.Party{
					Name: "counterparty 3",
					BIC:  "bic3",
					IBAN: "iban3",
				},
			},
		},
	}

	err = transferManager.ProcessTransfers(ctx, &request)
	assert.Error(t, err)

	qontoAccountAfterProcessing, err := mysqlStorage.FindAccount(ctx, qontoAccountID)
	require.NoError(t, err)
	assert.Equal(t, accountBalance, qontoAccountAfterProcessing.BalanceCents)

	var expectedBalance int64 = accountBalance
	assert.Equal(t, expectedBalance, qontoAccountAfterProcessing.BalanceCents)

	transactions, err := mysqlStorage.FindAccountTransactions(ctx, qontoAccountID)
	require.NoError(t, err)
	assert.Len(t, transactions, 0)
}
