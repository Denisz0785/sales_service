package web

import (
	"log"
	"net/http"
	"sales_service/internal/platform/web"
)

func Errors(log *log.Logger) web.Middleware {

	f := func(before web.Handler) web.Handler {

		h := func(w http.ResponseWriter, r *http.Request) error {

			if err := before(w, r); err != nil {
				log.Printf("Error: %v", err)

				if err := web.RespondError(w, err); err != nil {
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
