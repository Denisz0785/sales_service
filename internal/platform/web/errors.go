package web

import "github.com/pkg/errors"

type ErrorResponse struct {
	Error  string       `json:"error"`
	Fields []FieldError `json:"fields,omitempty"`
}

type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type Error struct {
	Status int
	Err    error
	Fields []FieldError
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func NewRequestError(err error, status int) error {
	return &Error{Err: err, Status: status}
}

// shutdown represents an error indicating that a shutdown was initiated.
type shutdown struct {
	// Message is the error message.
	Message string
}

// Error returns the error message.
func (s *shutdown) Error() string {
	return s.Message
}

// NewShutdownError returns a new shutdown error with the provided message.
func NewShutdownError(message string) error {
	return &shutdown{Message: message}
}

// IsShutdown returns true if the provided error is a shutdown error.
func IsShutdown(err error) bool {
	// Check if the error is an instance of shutdown.
	if _, ok := errors.Cause(err).(*shutdown); ok {
		return true
	}
	return false
}
