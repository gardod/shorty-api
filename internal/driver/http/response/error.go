package response

type Error string

func (e Error) Error() string { return string(e) }

const (
	ErrNotFound         = Error("Resource could not be found")
	ErrMethodNotAllowed = Error("Method not supported by the resource")
	ErrParse            = Error("Unable to process the request due to invalid syntax")
	ErrValidation       = Error("Unable to process the request due to invalid data")
	ErrInternal         = Error("The server encountered an unexpected condition")
)
