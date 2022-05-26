package core

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrNotEnoughFunds = Error("not enough funds")
)
