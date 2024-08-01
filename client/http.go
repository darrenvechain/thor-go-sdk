package client

import (
	"errors"
	"io"
	"net/http"
)

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"states"`
}

var ErrNotFound = &HttpError{Code: 404, Status: "not found", Message: "resource not found"}

func (e *HttpError) Error() string {
	return e.Message
}

func newHttpError(resp *http.Response) *HttpError {
	message := resp.Status

	if resp.Body != nil {
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err == nil {
			message = string(body)
		}
	}

	return &HttpError{
		Code:    resp.StatusCode,
		Status:  resp.Status,
		Message: message,
	}
}

func (e *HttpError) Is(target error) bool {
	var t *HttpError
	ok := errors.As(target, &t)
	if !ok {
		return false
	}
	return e.Code == t.Code
}
