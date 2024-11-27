package client

import (
	"testing"
)

func TestNewAuthsignalAPIError(t *testing.T) {
	errorCode := "bad_request"
	errorDescription := "An error occurred"
	statusCode := 400

	apiError := NewAuthsignalAPIError(errorCode, errorDescription, statusCode)

	if apiError == nil {
		t.Fatal("Expected NewAuthsignalAPIError to return a non-nil error")
	}

	if apiError.ErrorCode != errorCode {
		t.Errorf("Expected ErrorCode to be '%s', got '%s'", errorCode, apiError.ErrorCode)
	}

	if apiError.ErrorDescription != errorDescription {
		t.Errorf("Expected ErrorDescription to be '%s', got '%s'", errorDescription, apiError.ErrorDescription)
	}

	if apiError.StatusCode != statusCode {
		t.Errorf("Expected StatusCode to be %d, got %d", statusCode, apiError.StatusCode)
	}
}

func TestAuthsignalAPIError_Error(t *testing.T) {
	statusCode := 404
	errorDescription := "Not Found"

	apiError := NewAuthsignalAPIError("ERR404", errorDescription, statusCode)
	expectedErrorMessage := "AuthsignalException: 404 - Not Found"

	if apiError.Error() != expectedErrorMessage {
		t.Errorf("Expected Error() to return '%s', got '%s'", expectedErrorMessage, apiError.Error())
	}
}
