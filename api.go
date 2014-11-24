package apiservice

import (
	"appengine"
	"encoding/json"
	"net/http"
	"strings"
)

type APIError struct {
	IsError bool   `json:"is_error"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewAPIError(code int, message string) *APIError {
	return &APIError{
		true,
		code,
		http.StatusText(code),
		message,
	}
}

func (ae APIError) send(w http.ResponseWriter) {
	w.WriteHeader(ae.Code)
	renderJson(w, ae)
}

func renderJson(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(payload)
	if err != nil {
		panic(err.Error)
	}
	w.Write(bytes)
}

func apiSearchHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		NewAPIError(http.StatusBadRequest, "No query string found").send(w)
		return
	}
	nodes := searchNodes(c, q, 0)
	renderJson(w, nodes)
}