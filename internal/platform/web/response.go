package web

import (
	"encoding/json"
	"net/http"

	"github.com/go-faster/errors"
)

func Respond(w http.ResponseWriter, val interface{}, statusCode int) error {
	data, err := json.Marshal(val)
	if err != nil {
		return errors.Wrap(err, "marshal")
	}

	w.Header().Set("content-type", "application/json;charset=utf-8")
	w.WriteHeader(statusCode)
	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "write to client")
	}
	return nil
}

func RespondError(w http.ResponseWriter, err error) error {
	if webErr, ok := err.(*Error); ok {
		resp := ErrorResponse{
			Error: webErr.Err.Error(),
		}
		return Respond(w, resp, webErr.Status)
	}
	resp := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}

	return Respond(w, resp, http.StatusInternalServerError)
}
