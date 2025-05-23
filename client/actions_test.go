package client

import (
	"os"
	"testing"
)

var (
	actionTestConfig = TestConfig{
		apiSecretKey: os.Getenv("AUTHSIGNAL_API_SECRET"),
		apiUrl:       os.Getenv("AUTHSIGNAL_API_URL"),
	}
)

func TestActions(t *testing.T) {
	client := NewAuthsignalClient(actionTestConfig.apiSecretKey, actionTestConfig.apiUrl)

	redirectToSettings := true
	trackInput := TrackRequest{
		UserId: "user123",
		Action: "go-sdk-test",
		Attributes: &TrackAttributes{
			RedirectUrl:        "http://localhost:3000",
			RedirectToSettings: &redirectToSettings,
			Email:              "not-a-real-email@authsignal.com",
			PhoneNumber:        "1234567890",
			IpAddress:          "127.0.0.1",
			UserAgent:          "Authsignal-Go-SDK-Tests/1.0",
			DeviceId:           "device123",
			Custom: map[string]interface{}{
				"hello": "world",
			},
		},
	}

	trackResponse, err := client.Track(trackInput)
	if err != nil {
		t.Fatalf("Track failed: %v", err)
	}

	if trackResponse.State != "CHALLENGE_REQUIRED" {
		t.Errorf("Expected State to be 'CHALLENGE_REQUIRED', got %s", trackResponse.State)
	}

	if trackResponse.IdempotencyKey == "" {
		t.Error("Expected IdempotencyKey to be set, got empty string")
	}

	if trackResponse.IsEnrolled == nil || !*trackResponse.IsEnrolled {
		t.Error("Expected IsEnrolled to be true, got false")
	}

	if trackResponse.Url == "" {
		t.Error("Expected URL to be set, got empty string")
	}

	if trackResponse.Token == "" {
		t.Error("Expected Token to be set, got empty string")
	}

	if len(trackResponse.AllowedVerificationMethods) == 0 {
		t.Error("Expected AllowedVerificationMethods to be non-empty")
	}

	getActionInput := GetActionRequest{
		UserId:         "user123",
		Action:         "go-sdk-test",
		IdempotencyKey: trackResponse.IdempotencyKey,
	}

	getActionResponse, err := client.GetAction(getActionInput)
	if err != nil {
		t.Fatalf("GetAction failed: %v", err)
	}

	if getActionResponse.State != "CHALLENGE_REQUIRED" {
		t.Errorf("Expected State to be 'CHALLENGE_REQUIRED', got %s", getActionResponse.State)
	}

	if getActionResponse.Output == nil {
		t.Error("Expected Output to be set, got nil")
	}

	updateActionInput := UpdateActionRequest{
		UserId:         "user123",
		Action:         "go-sdk-test",
		IdempotencyKey: trackResponse.IdempotencyKey,
		Attributes: &ActionAttributes{
			State: "BLOCK",
		},
	}

	updateActionResponse, err := client.UpdateAction(updateActionInput)
	if err != nil {
		t.Fatalf("UpdateAction failed: %v", err)
	}
	if updateActionResponse.State != "BLOCK" {
		t.Errorf("Expected State to be 'BLOCK', got %s", updateActionResponse.State)
	}
}

func TestValidateChallenge(t *testing.T) {
	client := NewAuthsignalClient(actionTestConfig.apiSecretKey, actionTestConfig.apiUrl)

	redirectToSettings := true
	trackInput := TrackRequest{
		UserId: "user123",
		Action: "go-sdk-test",
		Attributes: &TrackAttributes{
			RedirectUrl:        "http://localhost:3000",
			RedirectToSettings: &redirectToSettings,
			Email:              "not-a-real-email@authsignal.com",
			PhoneNumber:        "1234567890",
			IpAddress:          "127.0.0.1",
			UserAgent:          "Authsignal-Go-SDK-Tests/1.0",
			DeviceId:           "device123",
		},
	}

	trackResponse, err := client.Track(trackInput)
	if err != nil {
		t.Fatalf("Track failed: %v", err)
	}

	validateInput := ValidateChallengeRequest{
		Token:  trackResponse.Token,
		UserId: trackInput.UserId,
		Action: trackInput.Action,
	}

	validateResponse, err := client.ValidateChallenge(validateInput)
	if err != nil {
		t.Fatalf("ValidateChallenge failed: %v", err)
	}

	if validateResponse.IsValid == nil {
		t.Error("Expected IsValid to be set, got nil")
	}

	if validateResponse.State != "CHALLENGE_REQUIRED" {
		t.Errorf("Expected State to be 'CHALLENGE_REQUIRED', got %s", validateResponse.State)
	}

	if validateResponse.UserId != trackInput.UserId {
		t.Errorf("Expected UserId to be '%s', got '%s'", trackInput.UserId, validateResponse.UserId)
	}

	if validateResponse.Action != trackInput.Action {
		t.Errorf("Expected Action to be '%s', got '%s'", trackInput.Action, validateResponse.Action)
	}

	if validateResponse.IdempotencyKey == "" {
		t.Error("Expected IdempotencyKey to be set, got empty string")
	}
}

func TestGetActionWithBadSecret(t *testing.T) {
	badConfig := TestConfig{
		apiSecretKey: "bad-secret",
		apiUrl:       actionTestConfig.apiUrl,
	}
	client := NewAuthsignalClient(badConfig.apiSecretKey, badConfig.apiUrl)

	input := GetActionRequest{
		UserId: "test-user",
		Action: "go-sdk-test",
	}

	_, err := client.GetAction(input)
	if err == nil {
		t.Fatal("Expected error with bad secret, got nil")
	}

	// Cast to AuthsignalAPIError to check specific error fields
	apiErr, ok := err.(*AuthsignalAPIError)
	if !ok {
		t.Fatalf("Expected error to be of type *AuthsignalAPIError, got %T", err)
	}

	if apiErr.StatusCode != 401 {
		t.Errorf("Expected status code 401, got %d", apiErr.StatusCode)
	}

	if apiErr.ErrorDescription != "The request is unauthorized. Check that your API key and region base URL are correctly configured." {
		t.Errorf("Expected error description 'The request is unauthorized. Check that your API key and region base URL are correctly configured.', got '%s'", apiErr.ErrorDescription)
	}

	if apiErr.ErrorCode != "unauthorized" {
		t.Errorf("Expected error code 'unauthorized', got '%s'", apiErr.ErrorCode)
	}

	if apiErr.Error() != "AuthsignalException: 401 - The request is unauthorized. Check that your API key and region base URL are correctly configured." {
		t.Errorf("Expected error string 'AuthsignalException: 401 - The request is unauthorized. Check that your API key and region base URL are correctly configured.', got '%s'", apiErr.Error())
	}
}

func TestTrackWithoutAttributes(t *testing.T) {
	client := NewAuthsignalClient(actionTestConfig.apiSecretKey, actionTestConfig.apiUrl)

	trackInput := TrackRequest{
		UserId: "user123",
		Action: "go-sdk-test",
	}

	trackResponse, err := client.Track(trackInput)
	if err != nil {
		t.Fatalf("Track failed: %v", err)
	}

	if trackResponse.State != "CHALLENGE_REQUIRED" {
		t.Errorf("Expected State to be 'CHALLENGE_REQUIRED', got %s", trackResponse.State)
	}

	if trackResponse.IdempotencyKey == "" {
		t.Error("Expected IdempotencyKey to be set, got empty string")
	}

	if trackResponse.Token == "" {
		t.Error("Expected Token to be set, got empty string")
	}
}
