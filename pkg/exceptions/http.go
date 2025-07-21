package exceptions

import "net/http"

type HTTPError interface {
	Error() string
	Status() int
}

type httpError struct {
	status  int
	message string
}

func (e *httpError) Error() string {
	return e.message
}

func (e *httpError) Status() int {
	return e.status
}

func NewHTTPError(status int, message string) HTTPError {
	return &httpError{
		status:  status,
		message: message,
	}
}

func NewBadRequest(message string) HTTPError {
	return NewHTTPError(http.StatusBadRequest, message)
}

func NewNotFound(message string) HTTPError {
	return NewHTTPError(http.StatusNotFound, message)
}

func NewInternalServerError(message string) HTTPError {
	return NewHTTPError(http.StatusInternalServerError, message)
}
