package core

type (
	TransferManager interface {
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
