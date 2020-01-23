package response

import (
	"encoding/json"
	"net/http"
)

type JSON struct {
	w http.ResponseWriter

	cookies []*http.Cookie
	code    int

	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

func NewJSON(w http.ResponseWriter) *JSON {
	return &JSON{w: w}
}

func (r *JSON) Send() error {
	r.w.WriteHeader(r.code)

	r.w.Header().Add("Content-Type", "application/json")
	for _, c := range r.cookies {
		http.SetCookie(r.w, c)
	}

	r.Status = http.StatusText(r.code)

	return json.NewEncoder(r.w).Encode(r)
}
