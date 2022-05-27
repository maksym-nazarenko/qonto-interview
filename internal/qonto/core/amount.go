package core

import (
	"errors"
	"strconv"
	"strings"
)

func ParseAmount(s string) (Amount, error) {
	if _, err := strconv.ParseFloat(s, 64); err != nil {
		return Amount{}, err
	}
	parts := strings.Split(s, ".")
	euros, cents := parts[0], ""
	if len(parts) == 2 {
		cents = parts[1]
	}
	if len(cents) > 2 {
		return Amount{}, errors.New("amount must contain at most 2 decimals after period")
	}
	eurosInt, err := strconv.Atoi(euros)
	if err != nil {
		return Amount{}, err
	}
	cents = cents + strings.Repeat("0", 2-len(cents))
	centsInt := 0
	centsInt, err = strconv.Atoi(cents)
	if err != nil {
		return Amount{}, err
	}

	return Amount{Cents: int64(eurosInt)*100 + int64(centsInt)}, nil
}

func (a *Amount) MarshalJSON() ([]byte, error) {
	euroes, cents := a.Cents/100, a.Cents%100
	result := strconv.FormatInt(euroes, 10)
	if cents > 0 {
		centsStr := strconv.FormatInt(cents, 10)
		result += "." + strings.Repeat("0", 2-len(centsStr)) + strings.TrimRight(centsStr, "0")
	}
	return []byte(result), nil
}

func (a *Amount) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	amount, err := ParseAmount(s)
	if err != nil {
		return err
	}
	*a = amount

	return nil
}
