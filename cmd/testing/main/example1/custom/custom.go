package custom

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"mime"
	"net/http"
)

func Only(method string, path string, h http.HandlerFunc) (string, http.HandlerFunc) {
	return path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			code := http.StatusMethodNotAllowed
			http.Error(w, http.StatusText(code), code)
			return
		}
		h(w, r)
	}
}

func Get(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			code := http.StatusMethodNotAllowed
			http.Error(w, http.StatusText(code), code)
			return
		}
		h(w, r)
	}
}

func Post(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			code := http.StatusMethodNotAllowed
			http.Error(w, http.StatusText(code), code)
			return
		}
		h(w, r)
	}
}

func Put(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			code := http.StatusMethodNotAllowed
			http.Error(w, http.StatusText(code), code)
			return
		}
		h(w, r)
	}
}

func Delete(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			code := http.StatusMethodNotAllowed
			http.Error(w, http.StatusText(code), code)
			return
		}
		h(w, r)
	}
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func BadRequest(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func ReadRaw(w http.ResponseWriter, r *http.Request) ([]byte, string) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		BadRequest(w, r)
	}
	return data, http.DetectContentType(data)
}

func ReadJSON(w http.ResponseWriter, r *http.Request, ptr any) {
	err := json.NewDecoder(r.Body).Decode(ptr)
	if err != nil {
		BadRequest(w, r)
	}
}

func ReadXML(w http.ResponseWriter, r *http.Request, ptr any) {
	err := xml.NewDecoder(r.Body).Decode(ptr)
	if err != nil {
		BadRequest(w, r)
	}
}

func WriteRaw(w http.ResponseWriter, r *http.Request, code int, data []byte) {
	if data == nil {
		data = []byte(http.StatusText(code))
	}
	w.Header().Set("Content-Type", http.DetectContentType(data))
	w.WriteHeader(code)
	_, err := w.Write(data)
	if err != nil {
		NotImplemented(w, r)
	}
}

type JSON map[string]any

type ErrJSON struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Error  error  `json:"error"`
}

func NewErrorJSON(code int, err error) ErrJSON {
	return ErrJSON{
		Code:   code,
		Status: http.StatusText(code),
		Error:  err,
	}
}

func WriteJSON(w http.ResponseWriter, r *http.Request, code int, data any) {
	w.Header().Set("Content-Type", mime.TypeByExtension(".json"))
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		NotImplemented(w, r)
	}
}

func WriteErrorJSON(w http.ResponseWriter, r *http.Request, code int, err error) {
	w.Header().Set("Content-Type", mime.TypeByExtension(".json"))
	w.WriteHeader(code)
	e := json.NewEncoder(w).Encode(
		ErrJSON{
			Code:   code,
			Status: http.StatusText(code),
			Error:  err,
		},
	)
	if e != nil {
		NotImplemented(w, r)
	}
}

func WriteXML(w http.ResponseWriter, r *http.Request, code int, data any) {
	w.Header().Set("Content-Type", mime.TypeByExtension(".xml"))
	w.WriteHeader(code)
	err := xml.NewEncoder(w).Encode(data)
	if err != nil {
		NotImplemented(w, r)
	}
}
