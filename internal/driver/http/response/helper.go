package response

import (
	"encoding/json"
	"net/http"

	vld "github.com/go-ozzo/ozzo-validation/v4"
)

func Prepare(w http.ResponseWriter, v interface{}, code int) *Response {
	resp := NewResponse(w)
	resp.SetStatusCode(code)

	if err, ok := v.(error); ok {
		switch err := err.(type) {
		case vld.Errors:
			resp.Error = ErrValidation.Error()
			resp.Details = err

		case *json.UnmarshalTypeError:
			resp.Error = ErrParse.Error()
			resp.Details = err

		default:
			resp.Error = err.Error()
		}

	} else {
		resp.Data = v
	}

	return resp
}

func JSON(w http.ResponseWriter, v interface{}, code int) {
	Prepare(w, v, code).JSON()
}

func Gob(w http.ResponseWriter, v interface{}, code int) {
	Prepare(w, v, code).Gob()
}
