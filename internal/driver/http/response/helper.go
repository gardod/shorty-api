package response

import (
	"net/http"
)

func Prepare(w http.ResponseWriter, v interface{}, code int) *Response {
	resp := NewResponse(w)
	resp.SetStatusCode(code)

	if err, ok := v.(*Error); ok {
		resp.Error = err.Error()
		if details := err.Unwrap(); details != nil {
			resp.Details = details
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
