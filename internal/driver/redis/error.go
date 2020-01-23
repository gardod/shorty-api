package redis

type Error string

func (e Error) Error() string { return string(e) }

const (
	ErrNotFound = Error("key does not exist")
)
