package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/maxim-nazarenko/qonto-interview/internal/qonto/core"
)

func (qapi *qontoAPI) HandleTransfers(w http.ResponseWriter, r *http.Request) {
	var request Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handleErrors(w, r, fmt.Errorf("error decoding request: %w: %v", ErrMalformedInput, err))
		return
	}

	coreRequest := core.Request{
		Party: core.Party{
			Name: request.OrganizationName,
			BIC:  request.OrganizationBIC,
			IBAN: request.OrganizationIBAN,
		},
		CreditTransfers: make([]core.Transfer, 0, len(request.CreditTransfers)),
	}

	for _, transfer := range request.CreditTransfers {
		if transfer.Currency != string(core.CURRENCY_EURO) {
			handleErrors(w, r, core.ErrInvalidCurrency)
		}
		coreRequest.CreditTransfers = append(coreRequest.CreditTransfers,
			core.Transfer{
				Amount:   transfer.Amount,
				Currency: core.Currency(transfer.Currency),
				CounterParty: core.Party{
					Name: transfer.CounterpartyName,
					BIC:  transfer.CounterpartyBIC,
					IBAN: transfer.CounterpartyIBAN,
				},
				Description: transfer.Description,
			})
	}
	if err := qapi.manager.ProcessTransfers(r.Context(), &coreRequest); err != nil {
		handleErrors(w, r, err)
		return
	}

	RespondCode(w, r, http.StatusCreated, "operation succeeded")

}
