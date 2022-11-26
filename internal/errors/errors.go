package errors

import (
	"net/http"
)

// All possible errors withing the service

type errorCode string

var (
	Internal     errorCode = "internal"
	NotFound     errorCode = "not_found"
	InvalidInput errorCode = "invalid_input"
)

// ApiError encapsulates error data to be sent out of the service via HTTP
type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewApiError(code errorCode, err error) ApiError {
	apiError := ApiError{Message: err.Error()}
	switch code {
	case InvalidInput:
		apiError.Code = http.StatusBadRequest
	case NotFound:
		apiError.Code = http.StatusNotFound
	default:
		apiError.Code = http.StatusInternalServerError
	}

	return apiError
}
