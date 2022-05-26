package core

import (
	"encoding/json"
	"strconv"
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

func TestAmountMarshalJSON(t *testing.T) {
	cases := []struct {
		name           string
		input          Amount
		expectedResult string
		expectError    bool
	}{
		{
			name:           "euros only",
			input:          Amount{Cents: 1000},
			expectedResult: "10",
		},
		{
			name:           "euros and .02 cents",
			input:          Amount{Cents: 2002},
			expectedResult: "20.02",
		},
		{
			name:           "euros and .50 cents",
			input:          Amount{Cents: 2050},
			expectedResult: "20.5",
		},
		{
			name:           "euros and .99 cents",
			input:          Amount{Cents: 2099},
			expectedResult: "20.99",
		},
		{
			name:           "zero",
			input:          Amount{Cents: 0},
			expectedResult: "0",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := json.Marshal(&tc.input)
			if (err == nil) != (tc.expectError == false) {
				t.Errorf(`
				expected error to be %v, got %v
				`, tc.expectError, err)
				return
			}

			assert.Equal(t, tc.expectedResult, string(result))
		})
	}
}

func TestAmountUnmarshalJSON(t *testing.T) {
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
			name:        "fail, more than 2 digits of cents",
			input:       "30.123",
			expectError: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var result Amount
			err := json.Unmarshal([]byte(strconv.Quote(tc.input)), &result)
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
