package client

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

	if client.Retries != DefaultRetries {
		t.Errorf("Expected retries to be %d, got %d", DefaultRetries, client.Retries)
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

func withNoRetryDelay(t *testing.T) {
	t.Helper()
	original := retryBaseDelay
	retryBaseDelay = 0
	t.Cleanup(func() { retryBaseDelay = original })
}

func TestRetriesSafeRequestsTwiceOn5xx(t *testing.T) {
	withNoRetryDelay(t)
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprint(w, `{"errorCode":"unavailable"}`)
			return
		}
		fmt.Fprint(w, `{}`)
	}))
	defer server.Close()

	client := NewAuthsignalClient("secret", server.URL)
	if _, err := client.get("/users/user"); err != nil {
		t.Fatalf("Expected request to succeed, got %v", err)
	}
	if attempts != 3 {
		t.Fatalf("Expected 3 attempts, got %d", attempts)
	}
}

func TestRetries429Responses(t *testing.T) {
	withNoRetryDelay(t)
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		if attempts == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprint(w, `{"errorCode":"rate_limited"}`)
			return
		}
		fmt.Fprint(w, `{}`)
	}))
	defer server.Close()

	client := NewAuthsignalClient("secret", server.URL)
	if _, err := client.get("/users/user"); err != nil {
		t.Fatalf("Expected request to succeed, got %v", err)
	}
	if attempts != 2 {
		t.Fatalf("Expected 2 attempts, got %d", attempts)
	}
}

func TestRetriesTransientNetworkFailures(t *testing.T) {
	withNoRetryDelay(t)
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		if attempts == 1 {
			hijacker, ok := w.(http.Hijacker)
			if !ok {
				t.Fatal("test server does not support hijacking")
			}
			connection, _, err := hijacker.Hijack()
			if err != nil {
				t.Fatalf("failed to hijack connection: %v", err)
			}
			connection.Close()
			return
		}
		fmt.Fprint(w, `{}`)
	}))
	defer server.Close()

	client := NewAuthsignalClient("secret", server.URL)
	if _, err := client.get("/users/user"); err != nil {
		t.Fatalf("Expected request to succeed, got %v", err)
	}
	if attempts != 2 {
		t.Fatalf("Expected 2 attempts, got %d", attempts)
	}
}

func TestRetriesIdempotentWrites(t *testing.T) {
	withNoRetryDelay(t)
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		if attempts == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprint(w, `{"errorCode":"unavailable"}`)
			return
		}
		fmt.Fprint(w, `{}`)
	}))
	defer server.Close()

	client := NewAuthsignalClient("secret", server.URL)
	if _, err := client.post("/users/user/actions/withdrawal", bytes.NewBufferString(`{"idempotencyKey":"key"}`)); err != nil {
		t.Fatalf("Expected request to succeed, got %v", err)
	}
	if attempts != 2 {
		t.Fatalf("Expected 2 attempts, got %d", attempts)
	}
}

func TestDoesNotRetryNonIdempotentWritesOr499(t *testing.T) {
	withNoRetryDelay(t)
	for _, test := range []struct {
		name   string
		method string
		status int
		body   string
	}{
		{name: "write", method: http.MethodPost, status: http.StatusServiceUnavailable, body: `{}`},
		{name: "499", method: http.MethodGet, status: 499, body: ``},
	} {
		t.Run(test.name, func(t *testing.T) {
			attempts := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				attempts++
				w.WriteHeader(test.status)
				fmt.Fprint(w, `{"errorCode":"failed"}`)
			}))
			defer server.Close()

			client := NewAuthsignalClient("secret", server.URL)
			_, _ = client.makeRequest(test.method, "/test", bytes.NewBufferString(test.body))
			if attempts != 1 {
				t.Fatalf("Expected 1 attempt, got %d", attempts)
			}
		})
	}
}

func TestAllowsRetriesToBeDisabled(t *testing.T) {
	withNoRetryDelay(t)
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, `{"errorCode":"unavailable"}`)
	}))
	defer server.Close()

	client := NewAuthsignalClient("secret", server.URL)
	client.Retries = 0
	_, _ = client.get("/users/user")
	if attempts != 1 {
		t.Fatalf("Expected 1 attempt, got %d", attempts)
	}
}

func TestRetryDefaults(t *testing.T) {
	client := NewAuthsignalClient("secret", "")
	if client.Client.Timeout != 10*time.Second || client.Retries != 2 {
		t.Fatalf("Unexpected retry defaults: timeout=%v retries=%d", client.Client.Timeout, client.Retries)
	}
}
