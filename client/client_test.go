package client

import (
	"testing"
)

func TestNew(t *testing.T) {
	client := NewAuthsignalClient("secret", "https://api.authsignal.com")

	if client.ApiSecretKey != "secret" {
		t.Errorf("Expected apiSecretKey to be 'secret', got %s", client.ApiSecretKey)
	}

	if client.ApiUrl != "https://api.authsignal.com" {
		t.Errorf("Expected apiUrl to be 'https://api.authsignal.com', got %s", client.ApiUrl)
	}

	if client.Client == nil {
		t.Error("Expected http client to be initialized")
	}

	if client.Client.Timeout != RequestTimeout {
		t.Errorf("Expected timeout to be %v, got %v", RequestTimeout, client.Client.Timeout)
	}
}

func TestDefaultHeaders(t *testing.T) {
	client := NewAuthsignalClient("secret", "https://api.authsignal.com")
	headers := client.defaultHeaders()

	expectedHeaders := map[string][]string{
		"Accept":       {"*/*"},
		"Content-Type": {"application/json"},
		"User-Agent":   {"authsignalgo/v1"},
	}

	for key, expected := range expectedHeaders {
		if actual := headers[key]; len(actual) != 1 || actual[0] != expected[0] {
			t.Errorf("Expected header %s to be %v, got %v", key, expected, actual)
		}
	}
}
