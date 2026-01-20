package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var (
	userTestConfig = TestConfig{
		apiSecretKey: os.Getenv("AUTHSIGNAL_API_SECRET"),
		apiUrl:       os.Getenv("AUTHSIGNAL_API_URL"),
	}
)

// Unit tests for QueryUsers that don't require API credentials

func TestQueryUsersBuildsCorrectURLWithEmail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		if !strings.HasPrefix(r.URL.Path, "/users") {
			t.Errorf("Expected path to start with /users, got %s", r.URL.Path)
		}

		email := r.URL.Query().Get("email")
		if email != "test@example.com" {
			t.Errorf("Expected email query param 'test@example.com', got '%s'", email)
		}

		response := QueryUsersResponse{
			Users: []QueryUsersResponseUser{
				{UserId: "user-1", Email: "test@example.com", EmailVerified: true},
			},
			LastEvaluatedUserId: "user-1",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewAuthsignalClient("test-secret", server.URL)
	result, err := client.QueryUsers(QueryUsersRequest{Email: "test@example.com"})

	if err != nil {
		t.Fatalf("QueryUsers failed: %v", err)
	}

	if len(result.Users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(result.Users))
	}

	if result.Users[0].Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", result.Users[0].Email)
	}
}

func TestQueryUsersBuildsCorrectURLWithAllParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		// Verify all query parameters
		if query.Get("username") != "testuser" {
			t.Errorf("Expected username 'testuser', got '%s'", query.Get("username"))
		}
		if query.Get("email") != "test@example.com" {
			t.Errorf("Expected email 'test@example.com', got '%s'", query.Get("email"))
		}
		if query.Get("phoneNumber") != "+1234567890" {
			t.Errorf("Expected phoneNumber '+1234567890', got '%s'", query.Get("phoneNumber"))
		}
		if query.Get("token") != "some-token" {
			t.Errorf("Expected token 'some-token', got '%s'", query.Get("token"))
		}
		if query.Get("limit") != "10" {
			t.Errorf("Expected limit '10', got '%s'", query.Get("limit"))
		}
		if query.Get("lastEvaluatedUserId") != "prev-user" {
			t.Errorf("Expected lastEvaluatedUserId 'prev-user', got '%s'", query.Get("lastEvaluatedUserId"))
		}

		response := QueryUsersResponse{Users: []QueryUsersResponseUser{}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewAuthsignalClient("test-secret", server.URL)
	limit := 10
	_, err := client.QueryUsers(QueryUsersRequest{
		Username:            "testuser",
		Email:               "test@example.com",
		PhoneNumber:         "+1234567890",
		Token:               "some-token",
		Limit:               &limit,
		LastEvaluatedUserId: "prev-user",
	})

	if err != nil {
		t.Fatalf("QueryUsers failed: %v", err)
	}
}

func TestQueryUsersNoParamsNoQueryString(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Errorf("Expected no query string, got '%s'", r.URL.RawQuery)
		}

		response := QueryUsersResponse{Users: []QueryUsersResponseUser{}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewAuthsignalClient("test-secret", server.URL)
	_, err := client.QueryUsers(QueryUsersRequest{})

	if err != nil {
		t.Fatalf("QueryUsers failed: %v", err)
	}
}

func TestQueryUsersReturnsCorrectResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := QueryUsersResponse{
			Users: []QueryUsersResponseUser{
				{
					UserId:              "user-1",
					Email:               "user1@example.com",
					EmailVerified:       true,
					PhoneNumber:         "+1234567890",
					PhoneNumberVerified: false,
					Username:            "user1",
				},
				{
					UserId:              "user-2",
					Email:               "user2@example.com",
					EmailVerified:       false,
					PhoneNumber:         "",
					PhoneNumberVerified: false,
					Username:            "user2",
				},
			},
			LastEvaluatedUserId: "user-2",
			TokenPayload:        map[string]string{"sub": "user-2"},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewAuthsignalClient("test-secret", server.URL)
	result, err := client.QueryUsers(QueryUsersRequest{Email: "example.com"})

	if err != nil {
		t.Fatalf("QueryUsers failed: %v", err)
	}

	if len(result.Users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(result.Users))
	}

	if result.Users[0].UserId != "user-1" {
		t.Errorf("Expected first user ID 'user-1', got '%s'", result.Users[0].UserId)
	}

	if result.Users[0].EmailVerified != true {
		t.Errorf("Expected first user EmailVerified true")
	}

	if result.Users[1].PhoneNumberVerified != false {
		t.Errorf("Expected second user PhoneNumberVerified false")
	}

	if result.LastEvaluatedUserId != "user-2" {
		t.Errorf("Expected LastEvaluatedUserId 'user-2', got '%s'", result.LastEvaluatedUserId)
	}

	if result.TokenPayload == nil {
		t.Errorf("Expected TokenPayload to be non-nil")
	}
}

func TestQueryUsersEmptyResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := QueryUsersResponse{Users: []QueryUsersResponseUser{}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewAuthsignalClient("test-secret", server.URL)
	result, err := client.QueryUsers(QueryUsersRequest{Email: "nonexistent@example.com"})

	if err != nil {
		t.Fatalf("QueryUsers failed: %v", err)
	}

	if len(result.Users) != 0 {
		t.Errorf("Expected 0 users, got %d", len(result.Users))
	}
}

func TestQueryUsersPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		if query.Get("limit") != "5" {
			t.Errorf("Expected limit '5', got '%s'", query.Get("limit"))
		}
		if query.Get("lastEvaluatedUserId") != "user-1" {
			t.Errorf("Expected lastEvaluatedUserId 'user-1', got '%s'", query.Get("lastEvaluatedUserId"))
		}

		response := QueryUsersResponse{
			Users: []QueryUsersResponseUser{
				{UserId: "user-2", Email: "user2@example.com"},
			},
			LastEvaluatedUserId: "user-2",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewAuthsignalClient("test-secret", server.URL)
	limit := 5
	result, err := client.QueryUsers(QueryUsersRequest{
		Email:               "example.com",
		Limit:               &limit,
		LastEvaluatedUserId: "user-1",
	})

	if err != nil {
		t.Fatalf("QueryUsers failed: %v", err)
	}

	if result.LastEvaluatedUserId != "user-2" {
		t.Errorf("Expected LastEvaluatedUserId 'user-2', got '%s'", result.LastEvaluatedUserId)
	}
}

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

func TestQueryUsers(t *testing.T) {
	client := NewAuthsignalClient(userTestConfig.apiSecretKey, userTestConfig.apiUrl)

	// First create a user to query
	testEmail := "query-test-go@authsignal.com"
	updateUserInput := UpdateUserRequest{
		UserId: "query-test-user-go",
		Attributes: &UserAttributes{
			Email:       testEmail,
			DisplayName: "Query Test User Go",
		},
	}

	_, err := client.UpdateUser(updateUserInput)
	if err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}

	// Query by email (at least one of username, email, or phoneNumber is required)
	queryInput := QueryUsersRequest{
		Email: testEmail,
	}

	queryResponse, err := client.QueryUsers(queryInput)
	if err != nil {
		t.Fatalf("QueryUsers failed: %v", err)
	}

	if queryResponse.Users == nil {
		t.Fatalf("Expected Users to be non-nil")
	}

	if len(queryResponse.Users) == 0 {
		t.Fatalf("Expected at least one user in response")
	}

	// Verify the user we created is in the response
	found := false
	for _, user := range queryResponse.Users {
		if user.Email == testEmail {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected to find user with email '%s' in response", testEmail)
	}

	// Clean up
	deleteUserInput := DeleteUserRequest{
		UserId: "query-test-user-go",
	}

	err = client.DeleteUser(deleteUserInput)
	if err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}
}
