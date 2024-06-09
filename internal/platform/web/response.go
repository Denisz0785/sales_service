package web

import (
	"context"
	"encoding/json"

	"net/http"

	"github.com/pkg/errors"
)

// Respond writes the provided value to the http.ResponseWriter with the specified
// status code.
func Respond(ctx context.Context, w http.ResponseWriter, val interface{}, statusCode int) error {

	v, ok := ctx.Value(KeyValues).(*Values)

	if !ok {
		return errors.New("web value missing from context")
	}

	v.StatusCode = statusCode

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}
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
func RespondError(ctx context.Context, w http.ResponseWriter, err error) error {
	// Check if the error is of type *Error
	if webErr, ok := errors.Cause(err).(*Error); ok {
		// Create an ErrorResponse with the error message and fields
		er := ErrorResponse{
			Error:  webErr.Err.Error(),
			Fields: webErr.Fields,
		}
		// Write the error response to the client
		if err := Respond(ctx, w, er, webErr.Status); err != nil {
			return err
		}

		return nil
	}
	// Create a generic internal server error message
	er := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}

	// Write the error response to the client
	if err := Respond(ctx, w, er, http.StatusInternalServerError); err != nil {
		return err
	}
	return nil
}


