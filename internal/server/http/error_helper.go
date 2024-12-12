package http

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/willbicks/epigram/internal/logutils"
)

// serverError writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
// Logs error to logger.Error.
func (s *QuoteServer) serverError(w http.ResponseWriter, r *http.Request, err error) {
	s.Logger.ErrorContext(r.Context(), "internal server error", logutils.Error(err), "trace", debug.Stack())

	status := http.StatusInternalServerError
	http.Error(w, fmt.Sprintf("Error %v: %s", status, http.StatusText(status)), status)
}

// clientError is a helper which sends a specific status code and corresponding
// description to the user, as well as the string representation of the optional
// err parameter. Logs error to logger.Debug.
func (s *QuoteServer) clientError(w http.ResponseWriter, r *http.Request, err error, code int) {
	var status string

	if err != nil {
		status = fmt.Sprintf("%v: %s", code, err.Error())
	} else {
		status = fmt.Sprintf("Error %v: %v", code, http.StatusText(code))
	}

	s.Logger.DebugContext(r.Context(), status)
	http.Error(w, status, code)
}

// notFound is a helper that wraps clientError and writes a 404 not found error.
func (s *QuoteServer) notFoundError(w http.ResponseWriter, r *http.Request) {
	s.clientError(w, r, nil, http.StatusNotFound)
}

// unsupportedMethod is a helper that wraps clientError and writes a 405 method not allowed error.
func (s *QuoteServer) methodNotAllowedError(w http.ResponseWriter, r *http.Request) {
	s.clientError(w, r, nil, http.StatusMethodNotAllowed)
}
