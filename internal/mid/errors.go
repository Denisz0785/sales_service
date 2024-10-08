package mid

import (
	"context"
	"log"
	"net/http"
	"sales_service/internal/platform/web"

	"go.opencensus.io/trace"
)

// Errors is a middleware function that wraps the provided web.Handler and
// logs and responds to errors that occur during the execution of the handler.
func Errors(log *log.Logger) web.Middleware {

	// Error handling middleware function
	f := func(before web.Handler) web.Handler {

		// New handler function that wraps the input web.Handler and handles
		// any errors that occur during its execution.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx, span := trace.StartSpan(ctx, "internal.mid.Errors")
			defer span.End()

			if err := before(ctx, w, r); err != nil {
				log.Printf("Error: %+v", err)

				// If there is an error, call the web.RespondError function
				// to create an error response and write it to the client.
				if err := web.RespondError(ctx, w, err); err != nil {
					return err
				}
				if web.IsShutdown(err) {
					return err
				}
			}
			return nil
		}
		return h
	}
	return f

}

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
