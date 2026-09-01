package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DEFAULT_API_URL = "https://api.authsignal.com/v1"
const RequestTimeout = 10 * time.Second
const ConnectTimeout = 3 * time.Second
const DefaultRetries = 2

var retryBaseDelay = 100 * time.Millisecond

type Client struct {
	ApiSecretKey string
	ApiUrl       string
	Client       *http.Client
	Retries      int
}

func NewAuthsignalClient(apiSecretKey string, apiUrl string) Client {
	if apiUrl == "" {
		apiUrl = DEFAULT_API_URL
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = (&net.Dialer{Timeout: ConnectTimeout, KeepAlive: 30 * time.Second}).DialContext

	return Client{
		ApiSecretKey: apiSecretKey,
		ApiUrl:       apiUrl,
		Client:       &http.Client{Timeout: RequestTimeout, Transport: transport},
		Retries:      DefaultRetries,
	}
}

func (c Client) defaultHeaders() http.Header {
	return http.Header{
		"Accept":       {"*/*"},
		"Content-Type": {"application/json"},
		"User-Agent":   {"authsignalgo/v1"},
	}
}

func (c Client) get(path string) ([]byte, error) {
	return c.makeRequest("GET", path, nil)
}

func (c Client) post(path string, body io.Reader) ([]byte, error) {
	return c.makeRequest("POST", path, body)
}

func (c Client) patch(path string, body io.Reader) ([]byte, error) {
	return c.makeRequest("PATCH", path, body)
}

func (c Client) delete(path string) ([]byte, error) {
	return c.makeRequest("DELETE", path, nil)
}

func (c Client) makeRequest(method, path string, body io.Reader) ([]byte, error) {
	var bodyBytes []byte
	var err error
	if body != nil {
		bodyBytes, err = io.ReadAll(body)
		if err != nil {
			return nil, err
		}
	}

	replayable := isReplayable(method, path, bodyBytes)
	for retryCount := 0; ; retryCount++ {
		req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.ApiUrl, path), bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, err
		}

		req.Header = c.defaultHeaders()
		req.SetBasicAuth(c.ApiSecretKey, "")

		resp, requestErr := c.Client.Do(req)
		if requestErr != nil {
			if replayable && retryCount < c.Retries {
				time.Sleep(retryDelay(retryCount, nil))
				continue
			}
			return nil, requestErr
		}

		responseBody, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			if replayable && retryCount < c.Retries {
				time.Sleep(retryDelay(retryCount, resp))
				continue
			}
			return nil, readErr
		}

		if replayable && retryCount < c.Retries && isRetryableStatus(resp.StatusCode) {
			time.Sleep(retryDelay(retryCount, resp))
			continue
		}

		if resp.StatusCode > 299 {
			var apiErr AuthsignalAPIError
			err := json.Unmarshal(responseBody, &apiErr)
			apiErr.StatusCode = resp.StatusCode

			if err != nil {
				return nil, err
			}

			return nil, &apiErr
		}

		return responseBody, nil
	}
}

func isReplayable(method, path string, body []byte) bool {
	method = strings.ToUpper(method)
	if method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions {
		return true
	}

	var payload map[string]interface{}
	if len(body) > 0 && json.Unmarshal(body, &payload) == nil {
		if key, ok := payload["idempotencyKey"].(string); ok && key != "" {
			return true
		}
	}

	parts := strings.Split(strings.Trim(path, "/"), "/")
	return method == http.MethodPatch && len(parts) >= 5 && parts[len(parts)-3] == "actions"
}

func isRetryableStatus(status int) bool {
	return status == http.StatusTooManyRequests || (status >= 500 && status <= 599)
}

func retryDelay(retryCount int, response *http.Response) time.Duration {
	baseDelay := retryBaseDelay * time.Duration(1<<retryCount)
	delay := baseDelay + time.Duration(rand.Int63n(maxInt64(1, int64(baseDelay/5))))

	if response != nil && response.StatusCode == http.StatusTooManyRequests {
		if retryAfter := response.Header.Get("Retry-After"); retryAfter != "" {
			if seconds, err := strconv.ParseFloat(retryAfter, 64); err == nil {
				retryAfterDelay := time.Duration(seconds * float64(time.Second))
				if retryAfterDelay > delay {
					delay = retryAfterDelay
				}
			} else if date, err := http.ParseTime(retryAfter); err == nil && time.Until(date) > delay {
				delay = time.Until(date)
			}
		}
	}

	return delay
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
