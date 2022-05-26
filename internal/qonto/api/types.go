package api

import (
	"github.com/maxim-nazarenko/qonto-interview/internal/qonto/core"
)

type (
	API interface {
	}

	qontoAPI struct {
		manager core.TransferManager
	}

	Transfer struct {
		Amount           core.Amount `json:"amount,omitempty"`
		Currency         string      `json:"currency,omitempty"`
		Description      string      `json:"description,omitempty"`
		CounterpartyName string      `json:"counterparty_name,omitempty"`
		CounterpartyBIC  string      `json:"counterparty_bic,omitempty"`
		CounterpartyIBAN string      `json:"counterparty_iban,omitempty"`
	}

	Request struct {
		OrganizationName string     `json:"organization_name,omitempty"`
		OrganizationBIC  string     `json:"organization_bic,omitempty"`
		OrganizationIBAN string     `json:"organization_iban,omitempty"`
		CreditTransfers  []Transfer `json:"credit_transfers,omitempty"`
	}
)
