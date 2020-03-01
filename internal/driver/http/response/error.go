package response

type Error struct {
	Text    string
	Details error
}

func (e *Error) Error() string { return e.Text }

func (e *Error) Unwrap() error { return e.Details }

func (e *Error) WithDetails(err error) *Error { return &Error{e.Text, err} }

var (
	ErrNotFound         = &Error{Text: "Resource could not be found"}
	ErrMethodNotAllowed = &Error{Text: "Method not supported by the resource"}
	ErrParse            = &Error{Text: "Unable to process the request due to invalid syntax"}
	ErrValidation       = &Error{Text: "Unable to process the request due to invalid data"}
	ErrInternal         = &Error{Text: "The server encountered an unexpected condition"}
)
