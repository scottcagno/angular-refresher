package rest

import (
	"encoding/json"
	"net/http"
)

type ContentType string

const (
	ContentJSON ContentType = "application/json"
	ContentXML  ContentType = "application/xml"
	ContentHTML ContentType = "text/html"
	ContentText ContentType = "text/plain"
)

type Payload interface {
	Marshal() ([]byte, error)
}

func WriteResponse(w http.ResponseWriter, code int, payload Payload) {
	data, err := payload.Marshal()
	if err != nil {
		code = http.StatusExpectationFailed
		w.WriteHeader(code)
		w.Write([]byte(http.StatusText(code)))
		return
	}
	w.WriteHeader(code)
	contentType := http.DetectContentType(data)
	w.Header().Set("Content-Type", contentType)
	w.Write(data)
	return
}

type Error struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewError(code int, msg string) *Error {
	return &Error{
		Code:    code,
		Status:  http.StatusText(code),
		Message: msg,
	}
}

func (e *Error) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

type Raw struct {
	Data any
}

func NewRaw(data any) *Raw {
	return &Raw{
		Data: data,
	}
}

func (raw *Raw) Marshal() ([]byte, error) {
	return json.Marshal(raw.Data)
}
