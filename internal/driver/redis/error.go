package redis

const (
	ErrNotFound = Error("key does not exist")
)

type Error string

func (e Error) Error() string { return string(e) }
