package repository

type Error string

func (e Error) Error() string { return string(e) }

const (
	ErrUniqueViolation = Error("unique constraint violation")
	ErrNoResults       = Error("no results found")
)
