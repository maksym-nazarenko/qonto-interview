package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/maxim-nazarenko/qonto-interview/internal/qonto/core"
	"github.com/stretchr/testify/assert"
)

func TestHandleTransfers(t *testing.T) {
	testCases := []struct {
		name           string
		body           string
		api            *qontoAPI
		expectedStatus int
	}{
		{
			name: "happy",
			api:  NewAPI(newMockManager()),
			body: `
			{
				"organization_name": "ACME Corp",
				"organization_bic": "OIVUSCLQXXX",
				"organization_iban": "FR10474608000002006107XXXXX",
				"credit_transfers": [
				  {
					"amount": "14.5",
					"currency": "EUR",
					"counterparty_name": "Bip Bip",
					"counterparty_bic": "CRLYFRPPTOU",
					"counterparty_iban": "EE383680981021245685",
					"description": "Wonderland/4410"
				  },
				  {
					"amount": "61238",
					"currency": "EUR",
					"counterparty_name": "Wile E Coyote",
					"counterparty_bic": "ZDRPLBQI",
					"counterparty_iban": "DE9935420810036209081725212",
					"description": "//TeslaMotors/Invoice/12"
				  },
				  {
					"amount": "999",
					"currency": "EUR",
					"counterparty_name": "Bugs Bunny",
					"counterparty_bic": "RNJZNTMC",
					"counterparty_iban": "FR0010009380540930414023042",
					"description": "2020 09 24/2020 09 25/GoldenCarrot/"
				  }
				]
			  }
			`,
			expectedStatus: http.StatusCreated,
		},
		{
			name: "manager returns not enough funds error",
			api:  NewAPI(newMockManager().WithError(core.ErrNotEnoughFunds)),
			body: `
			{
				"organization_name": "ACME Corp",
				"organization_bic": "OIVUSCLQXXX",
				"organization_iban": "FR10474608000002006107XXXXX",
				"credit_transfers": [
				  {
					"amount": "14.5",
					"currency": "EUR",
					"counterparty_name": "Bip Bip",
					"counterparty_bic": "CRLYFRPPTOU",
					"counterparty_iban": "EE383680981021245685",
					"description": "Wonderland/4410"
				  },
				  {
					"amount": "61238",
					"currency": "EUR",
					"counterparty_name": "Wile E Coyote",
					"counterparty_bic": "ZDRPLBQI",
					"counterparty_iban": "DE9935420810036209081725212",
					"description": "//TeslaMotors/Invoice/12"
				  },
				  {
					"amount": "999",
					"currency": "EUR",
					"counterparty_name": "Bugs Bunny",
					"counterparty_bic": "RNJZNTMC",
					"counterparty_iban": "FR0010009380540930414023042",
					"description": "2020 09 24/2020 09 25/GoldenCarrot/"
				  }
				]
			  }
			`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "empty input is not valid",
			api:  NewAPI(newMockManager()),
			body: `
			`,
			expectedStatus: http.StatusBadRequest,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "http://localhost", strings.NewReader(tc.body))
			tc.api.HandleTransfers(w, r)

			if !assert.Equal(t, tc.expectedStatus, w.Result().StatusCode) {
				body, _ := ioutil.ReadAll(w.Result().Body)
				t.Error(string(body))
			}
		})
	}
}
