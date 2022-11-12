package rest

import (
	"fmt"
	"net/http"
)

type Route struct {
	Method  string
	Pattern string
	Handler http.Handler
}

type routeEntry struct {
	method  string
	pattern string
	handler http.Handler
}

func (m routeEntry) String() string {
	if m.method == http.MethodGet {
		return fmt.Sprintf("[%s]&nbsp;&nbsp;&nbsp;&nbsp;<a href=\"%s\">%s</a>", m.method, m.pattern, m.pattern)
	}
	if m.method == http.MethodPost {
		return fmt.Sprintf("[%s]&nbsp;&nbsp;&nbsp;%s", m.method, m.pattern)
	}
	if m.method == http.MethodPut {
		return fmt.Sprintf("[%s]&nbsp;&nbsp;&nbsp;&nbsp;%s", m.method, m.pattern)
	}
	if m.method == http.MethodDelete {
		return fmt.Sprintf("[%s]&nbsp;%s", m.method, m.pattern)
	}
	return fmt.Sprintf("[%s]&nbsp;%s", m.method, m.pattern)
}
