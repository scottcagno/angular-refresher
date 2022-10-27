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

func AsJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
