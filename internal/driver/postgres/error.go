package postgres

const (
	ErrTxActive   = Error("another transaction is already active")
	ErrTxInactive = Error("there is no active transaction")
)

type Error string

func (e Error) Error() string { return string(e) }
