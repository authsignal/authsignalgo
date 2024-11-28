package client

import (
	"os"
	"testing"
)

var (
	authenticatorTestConfig = TestConfig{
		apiSecretKey: os.Getenv("AUTHSIGNAL_API_SECRET"),
		apiUrl:       os.Getenv("AUTHSIGNAL_API_URL"),
	}
)

func TestAuthenticators(t *testing.T) {
	client := NewAuthsignalClient(authenticatorTestConfig.apiSecretKey, authenticatorTestConfig.apiUrl)

	enrollInput := EnrollVerifiedAuthenticatorRequest{
		UserId: "user12345",
		Attributes: &EnrollVerifiedAuthenticatorAttributes{
			VerificationMethod: "EMAIL_OTP",
			Email:              "not-a-real-email@authsignal.com",
		},
	}

	enrollResponse, err := client.EnrollVerifiedAuthenticator(enrollInput)
	if err != nil {
		t.Fatalf("EnrollVerifiedAuthenticator failed: %v", err)
	}

	if enrollResponse.Authenticator.UserAuthenticatorId == "" {
		t.Error("Expected UserAuthenticatorId to be set, got empty string")
	}
	if enrollResponse.Authenticator.VerificationMethod != enrollInput.Attributes.VerificationMethod {
		t.Errorf("Expected VerificationMethod to be '%s', got '%s'", enrollInput.Attributes.VerificationMethod, enrollResponse.Authenticator.VerificationMethod)
	}

	if enrollResponse.Authenticator.Email != enrollInput.Attributes.Email {
		t.Errorf("Expected Email to be '%s', got '%s'", enrollInput.Attributes.Email, enrollResponse.Authenticator.Email)
	}

	if enrollResponse.Authenticator.UserId != enrollInput.UserId {
		t.Errorf("Expected UserId to be '%s', got '%s'", enrollInput.UserId, enrollResponse.Authenticator.UserId)
	}

	getAuthInput := GetAuthenticatorsRequest{
		UserId: "user12345",
	}

	getAuthResponse, err := client.GetAuthenticators(getAuthInput)
	if err != nil {
		t.Fatalf("GetAuthenticators failed: %v", err)
	}

	if len(getAuthResponse) == 0 {
		t.Error("Expected Authenticators to be non-empty")
	}

	authenticator := getAuthResponse[0]
	if authenticator.UserAuthenticatorId == "" {
		t.Error("Expected UserAuthenticatorId to be set, got empty string")
	}

	if authenticator.VerificationMethod != enrollInput.Attributes.VerificationMethod {
		t.Errorf("Expected VerificationMethod to be '%s', got '%s'", enrollInput.Attributes.VerificationMethod, authenticator.VerificationMethod)
	}

	if authenticator.Email != enrollInput.Attributes.Email {
		t.Errorf("Expected Email to be '%s', got '%s'", enrollInput.Attributes.Email, authenticator.Email)
	}

	deleteInput := DeleteAuthenticatorRequest{
		UserId:              "user12345",
		UserAuthenticatorId: authenticator.UserAuthenticatorId,
	}

	err = client.DeleteAuthenticator(deleteInput)
	if err != nil {
		t.Fatalf("DeleteAuthenticator failed: %v", err)
	}

	getAuthResponseAfterDelete, err := client.GetAuthenticators(getAuthInput)
	if err != nil {
		t.Fatalf("GetAuthenticators after delete failed: %v", err)
	}

	for _, auth := range getAuthResponseAfterDelete {
		if auth.UserAuthenticatorId == authenticator.UserAuthenticatorId {
			t.Error("Expected authenticator to be deleted but it still exists")
		}
	}
}
