package http

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// serverError writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
// Logs error to logger.Warn.
func (s *CharismsServer) serverError(w http.ResponseWriter, r *http.Request, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	s.Logger.Warn(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError is a helper which sends a specific status code and corresponding
// description to the user, as well as the string representation of the optional
// err parameter. Logs error to logger.Debug.
func (s *CharismsServer) clientError(w http.ResponseWriter, r *http.Request, err error, code int) {
	var status string

	if err != nil {
		status = fmt.Sprintf("%s: %s", http.StatusText(code), err.Error())
	} else {
		status = http.StatusText(code)
	}

	s.Logger.Debug(status)
	http.Error(w, status, code)
}

// notFound is a helper that wraps clientError and writes a 404 not found error.
func (s *CharismsServer) notFoundError(w http.ResponseWriter, r *http.Request) {
	s.clientError(w, r, nil, http.StatusNotFound)
}

// unsupportedMethod is a helper that wraps clientError and writes a 405 method not allowed error.
func (s *CharismsServer) methodNotAllowedError(w http.ResponseWriter, r *http.Request) {
	s.clientError(w, r, nil, http.StatusMethodNotAllowed)
}
