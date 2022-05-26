package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAmount(t *testing.T) {
	cases := []struct {
		name           string
		input          string
		expectedResult Amount
		expectError    bool
	}{
		{
			name:           "euros only",
			input:          "10",
			expectedResult: Amount{Cents: 1000},
		},
		{
			name:           "euros and .00 cents",
			input:          "20.00",
			expectedResult: Amount{Cents: 2000},
		},
		{
			name:           "euros and .02 cents",
			input:          "20.02",
			expectedResult: Amount{Cents: 2002},
		},
		{
			name:           "euros and .5 cents",
			input:          "20.5",
			expectedResult: Amount{Cents: 2050},
		},
		{
			name:           "euros and .99 cents",
			input:          "20.99",
			expectedResult: Amount{Cents: 2099},
		},
		{
			name:           "euros and . cents produces 00 cents",
			input:          "30.",
			expectError:    false,
			expectedResult: Amount{Cents: 3000},
		},
		{
			name:        "fail, more than 2 digits of cents",
			input:       "30.123",
			expectError: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseAmount(tc.input)
			if (err == nil) != (tc.expectError == false) {
				t.Errorf(`
				expected error to be %v, got %v
				`, tc.expectError, err)
				return
			}
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
