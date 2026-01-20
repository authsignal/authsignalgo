package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const DEFAULT_API_URL = "https://api.authsignal.com/v1"
const RequestTimeout = 10 * time.Second

type Client struct {
	ApiSecretKey string
	ApiUrl       string
	Client       *http.Client
}

func NewAuthsignalClient(apiSecretKey string, apiUrl string) Client {
	if apiUrl == "" {
		apiUrl = DEFAULT_API_URL
	}
	return Client{ApiSecretKey: apiSecretKey, ApiUrl: apiUrl, Client: &http.Client{Timeout: RequestTimeout}}
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
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.ApiUrl, path), body)
	if err != nil {
		return nil, err
	}

	req.Header = c.defaultHeaders()
	req.SetBasicAuth(c.ApiSecretKey, "")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
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
