package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const RequestTimeout = 10 * time.Second

type Client struct {
	secret string
	apiUrl string
}

func New(secret, apiUrl string) Client {
	return Client{secret: secret, apiUrl: apiUrl}
}

func (c Client) defaultHeaders() http.Header {
	return http.Header{
		"Accept":       {"*/*"},
		"Content-Type": {"application/json"},
		"User-Agent":   {"authsignalgo/v1"},
	}
}

func (c Client) GetUser(request UserRequest) (UserResponse, error) {
	path := fmt.Sprintf("/users/%s", request.UserId)
	response, err := c.get(path)
	if err != nil {
		return UserResponse{}, err
	}

	var data UserResponse
	err2 := json.Unmarshal(response, &data)
	if err2 != nil {
		return UserResponse{}, err2
	}

	return data, nil
}

func (c Client) TrackAction(request TrackRequest) (TrackResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return TrackResponse{}, err
	}

	path := fmt.Sprintf("/users/%s/actions/%s", request.UserId, request.Action)
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
	path := fmt.Sprintf("/users/%s/actions/%s/%s", request.UserId, request.Action, request.IdempotencyKey)
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

	path := fmt.Sprintf("/users/%s/authenticators", request.UserId)
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

func (c Client) ValidateChallenge(request ValidateChallengeRequest) (ValidateChallengeResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return ValidateChallengeResponse{IsValid: false}, err
	}

	response, err2 := c.post("/validate", bytes.NewBuffer(body))
	if err2 != nil {
		return ValidateChallengeResponse{IsValid: false}, err2
	}

	var data ValidateChallengeResponse
	err3 := json.Unmarshal(response, &data)
	if err3 != nil {
		return ValidateChallengeResponse{IsValid: false}, err3
	}

	return data, nil
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
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.apiUrl, path), body)
	if err != nil {
		return nil, err
	}

	req.Header = c.defaultHeaders()
	req.SetBasicAuth(c.secret, "")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			return
		}
	}()

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("bad request to %s, http status code of %d, status was: %s", req.URL, resp.StatusCode, resp.Status)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
