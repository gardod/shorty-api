package response

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v3"
)

func SendSuccessResponse(w http.ResponseWriter, v interface{}, code int) error {
	resp := NewJSON(w)
	resp.code = code
	resp.Data = v

	return resp.Send()
}

func SendErrorResponse(w http.ResponseWriter, err error, code int) error {
	resp := NewJSON(w)
	resp.code = code

	switch err := err.(type) {
	case validation.Errors:
		resp.Error = ErrValidation.Error()
		resp.Details = err

	case *json.UnmarshalTypeError:
		resp.Error = ErrParse.Error()
		resp.Details = err

	default:
		resp.Error = err.Error()
	}

	return resp.Send()
}
