package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

/*
This is the AuthSignal Go Lang SDK.
The module wraps the AuthSignal APIs with a Go Lang implementation allowing easier integration.

Notable decisions:
- We have a 10-second timeout set, as response speed is controlled by Authsignal.
- We do not remap errors from any libraries are dependent on, changing such libraries should be considered
  a breaking change as users for the SDK are bound to the libraries versions.
  This was made because we do not have many libraries, and they are core golang libraries.
*/
// Todo it could be better to use Context for HTTP requests to put the Timeout + other config in.
// Todo deal with HTTP status code.

const RequestTimeout = 10 * time.Second

type Client struct {
	apiKey      string
	apiUrl      string
	redirectUrl string
}

func New(apiUrl, apiKey string, redirectUrl string) Client {
	return Client{apiKey: apiKey, apiUrl: apiUrl, redirectUrl: redirectUrl}
}

func (c Client) defaultHeaders() http.Header {
	return http.Header{
		"Accept":       {"*/*"},
		"Content-Type": {"application/json"},
		"User-Agent":   {c.userAgent()},
	}
}

func (c Client) userAgent() string {
	return "Authsignal Go v1" // todo make module version dynamic
}

func (c Client) GetUser(userId string) (string, error) {
	response, err := c.get(userId)
	return string(response), err
}

func (c Client) TrackAction(request TrackRequest) (TrackResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return TrackResponse{}, err
	}

	path := fmt.Sprintf("%s/actions/%s", request.UserId, request.Action)
	response, err2 := c.post(path, bytes.NewBuffer(body))
	if err2 != nil {
		return TrackResponse{}, err2
	}

	var data TrackResponse
	err3 := json.Unmarshal(response, &data)
	if err3 != nil {
		return TrackResponse{}, err3
	}

	return data, nil
}

func (c Client) GetAction(request GetActionRequest) (GetActionResponse, error) {
	path := fmt.Sprintf("%s/actions/%s/%s", request.UserId, request.Action, request.IdempotencyKey)
	response, err := c.get(path)
	if err != nil {
		return GetActionResponse{}, err
	}

	var data GetActionResponse
	err2 := json.Unmarshal(response, &data)
	if err2 != nil {
		return GetActionResponse{}, err2
	}

	return data, nil
}

func (c Client) EnrollVerifiedAuthenticator(request EnrollVerifiedAuthenticatorRequest) (EnrollVerifiedAuthenticatorResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return EnrollVerifiedAuthenticatorResponse{}, err
	}

	path := fmt.Sprintf("%s/authenticators", request.UserId)
	response, err2 := c.post(path, bytes.NewBuffer(body))
	if err2 != nil {
		return EnrollVerifiedAuthenticatorResponse{}, err2
	}

	var data EnrollVerifiedAuthenticatorResponse
	err3 := json.Unmarshal(response, &data)
	if err3 != nil {
		return EnrollVerifiedAuthenticatorResponse{}, err3
	}

	return data, nil
}

func (c Client) LoginWithEmail(request LoginWithEmailRequest) (LoginWithEmailResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return LoginWithEmailResponse{}, err
	}

	response, err2 := c.post(fmt.Sprintf("/email/%s/challenge", request.Email), bytes.NewBuffer(body))
	if err2 != nil {
		return LoginWithEmailResponse{}, err2
	}

	var data LoginWithEmailResponse
	err3 := json.Unmarshal(response, &data)
	if err3 != nil {
		return LoginWithEmailResponse{}, err3
	}

	return data, err
}

func (c Client) ValidateChallenge(request ValidateChallengeRequest) string {
	return "NOT YET IMPLEMENTED"
}

func (c Client) get(path string) ([]byte, error) {
	return c.makeRequest("GET", path, nil)
}

func (c Client) post(path string, body io.Reader) ([]byte, error) {
	return c.makeRequest("POST", path, body)
}

func (c Client) makeRequest(method, path string, body io.Reader) ([]byte, error) {
	client := http.Client{}
	client.Timeout = RequestTimeout
	req, err := http.NewRequest(method, fmt.Sprintf("%s/v1/users/%s", c.apiUrl, path), body)
	if err != nil {
		return nil, err
	}

	req.Header = c.defaultHeaders()
	req.SetBasicAuth(c.apiKey, "")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			return //todo handle error or pass up.
		}
	}()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
