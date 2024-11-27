package client

import "fmt"

type AuthsignalAPIError struct {
	ErrorCode        string `json:"errorCode"`
	ErrorDescription string `json:"errorDescription"`
	StatusCode       int    `json:"statusCode"`
}

func NewAuthsignalAPIError(errorCode string, errorDescription string, statusCode int) *AuthsignalAPIError {
	return &AuthsignalAPIError{
		ErrorCode:        errorCode,
		ErrorDescription: errorDescription,
		StatusCode:       statusCode,
	}
}

func (e *AuthsignalAPIError) Error() string {
	return fmt.Sprintf("AuthsignalException: %d - %s", e.StatusCode, e.ErrorDescription)
}
