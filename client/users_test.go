package client

import (
	"os"
	"testing"
)

var (
	userTestConfig = TestConfig{
		apiSecretKey: os.Getenv("AUTHSIGNAL_API_SECRET"),
		apiUrl:       os.Getenv("AUTHSIGNAL_API_URL"),
	}
)

func TestUsers(t *testing.T) {
	client := NewAuthsignalClient(userTestConfig.apiSecretKey, userTestConfig.apiUrl)

	updateUserInput := UpdateUserRequest{
		UserId: "a-new-user",
		Attributes: &UserAttributes{
			PhoneNumber: "9876543210",
			DisplayName: "A New User",
		},
	}

	updateUserResponse, err := client.UpdateUser(updateUserInput)
	if err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}

	if updateUserResponse.PhoneNumber != updateUserInput.Attributes.PhoneNumber {
		t.Errorf("Expected PhoneNumber to be '%s', got '%s'", updateUserInput.Attributes.PhoneNumber, updateUserResponse.PhoneNumber)
	}

	if updateUserResponse.DisplayName != updateUserInput.Attributes.DisplayName {
		t.Errorf("Expected DisplayName to be '%s', got '%s'", updateUserInput.Attributes.DisplayName, updateUserResponse.DisplayName)
	}

	getUserInput := GetUserRequest{
		UserId: "a-new-user",
	}

	getUserResponse, err := client.GetUser(getUserInput)
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}

	if getUserResponse.PhoneNumber != updateUserInput.Attributes.PhoneNumber {
		t.Errorf("Expected PhoneNumber to be '%s', got '%s'", updateUserInput.Attributes.PhoneNumber, getUserResponse.PhoneNumber)
	}

	if getUserResponse.DisplayName != updateUserInput.Attributes.DisplayName {
		t.Errorf("Expected DisplayName to be '%s', got '%s'", updateUserInput.Attributes.DisplayName, getUserResponse.DisplayName)
	}

	deleteUserInput := DeleteUserRequest{
		UserId: "a-new-user",
	}

	err = client.DeleteUser(deleteUserInput)
	if err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}

	getUserResponse, err = client.GetUser(getUserInput)
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}

	if getUserResponse.IsEnrolled == nil || *getUserResponse.IsEnrolled != false {
		t.Errorf("Expected IsEnrolled to be false after deletion")
	}
}
