package web

type ErrorResponse struct {
	Error string `json:"error"`
}

type Error struct {
	Status int
	Err    error
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func NewRequestError(err error, status int) error {
	return &Error{Err: err, Status: status}
}
