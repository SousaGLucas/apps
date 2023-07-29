package vos

import "errors"

var (
	ErrInvalidAccountEventType = errors.New("account event type not found")
)

var (
	AccountEventTypeCashIn  AccountEventType = "cash_in"
	AccountEventTypeCashOut AccountEventType = "cash_out"
)

type AccountEventType string

func ParseAccountEventType(name string) (AccountEventType, error) {
	switch AccountEventType(name) {
	case AccountEventTypeCashIn:
		return AccountEventTypeCashIn, nil
	case AccountEventTypeCashOut:
		return AccountEventTypeCashOut, nil
	default:
		return "", ErrInvalidAccountEventType
	}
}

func (t AccountEventType) String() string {
	return string(t)
}
