package mid

import (
	"errors"
	"log"
	"net/http"
	"sales_service/internal/platform/web"
	"time"
)

func Logger(log *log.Logger) web.Middleware {

	f := func(before web.Handler) web.Handler {

		h := func(w http.ResponseWriter, r *http.Request) error {
			v, ok := r.Context().Value(web.KeyValues).(*web.Values)

			if !ok {
				return errors.New("web value missing from context")
			}

			err := before(w, r)
			log.Printf("%d (%v) Method: %s  URL: %s", v.StatusCode, time.Since(v.Start), r.Method, r.URL.Path)

			return err
		}
		return h
	}
	return f

}

/*
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
*/
