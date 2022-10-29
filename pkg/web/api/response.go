package api

import (
	"encoding/json"
	"io"
	"net/http"
)

type Response struct {
	Code        int
	Message     string
	ContentType string
	Body        io.Writer
}

func ReadJSON(r *http.Request, ptr any) error {
	return json.NewDecoder(r.Body).Decode(ptr)
}

func WriteJSON(w http.ResponseWriter, code int, data any) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}
	return
}
