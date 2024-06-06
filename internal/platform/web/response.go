package web

import (
	"encoding/json"

	"net/http"

	"github.com/pkg/errors"
)

// Respond writes the provided value to the http.ResponseWriter with the specified
// status code.
func Respond(w http.ResponseWriter, val interface{}, statusCode int) error {
	// Marshal the value to JSON
	data, err := json.Marshal(val)
	if err != nil {
		return errors.Wrap(err, "marshal") // Wrap the error with a more descriptive message
	}

	// Set the content type header to application/json
	w.Header().Set("content-type", "application/json;charset=utf-8")

	// Set the status code of the response
	w.WriteHeader(statusCode)

	// Write the JSON data to the response
	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "write to client") // Wrap the error with a more descriptive message
	}

	return nil
}

// RespondError writes an error response to the client.
// If the error is of type *Error, the error message and fields are included in the response.
// Otherwise, a generic internal server error message is included.
func RespondError(w http.ResponseWriter, err error) error {
	// Check if the error is of type *Error
	if webErr, ok := errors.Cause(err).(*Error); ok {
		// Create an ErrorResponse with the error message and fields
		er := ErrorResponse{
			Error:  webErr.Err.Error(),
			Fields: webErr.Fields,
		}
		// Write the error response to the client
		if err := Respond(w, er, webErr.Status); err != nil {
			return err
		}

		return nil
	}
	// Create a generic internal server error message
	er := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}

	// Write the error response to the client
	if err := Respond(w, er, http.StatusInternalServerError); err != nil {
		return err
	}
	return nil
}

// ErrorResponse represents an error response that is sent to the client.
// type ErrorResponse struct {
// 	// Error is a string representation of the error.
// 	Error string `json:"error"`
// 	// Fields is a map of field names to error messages.
// 	Fields map[string]string `json:"fields,omitempty"`
// }
