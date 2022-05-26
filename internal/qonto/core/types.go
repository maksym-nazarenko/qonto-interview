package core

import "context"

type (
	TransferManager interface {
		ProcessTransfers(ctx context.Context, request *Request) error
	}

	Currency string

	Party struct {
		Name string
		BIC  string
		IBAN string
	}

	Transfer struct {
		Amount       Amount
		Currency     Currency
		Description  string
		CounterParty Party
	}

	Request struct {
		Party           Party
		CreditTransfers []Transfer
	}

	Amount struct {
		Cents int64
	}
)

const (
	CURRENCY_EURO Currency = "EUR"
)
