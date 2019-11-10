package api

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrEmptyArgument = errors.New("(api) argument cannot be empty")
)

type ErrorResponse struct {
	Response     *http.Response
	ErrorElement ErrorElement
}

type ErrorElement struct {
	Error string `json:"error"`
	Code  int    `json:"code,omitempty"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %+v",
		r.Response.Request.Method, sanitizeURL(r.Response.Request.URL), r.Response.StatusCode, r.ErrorElement)
}
