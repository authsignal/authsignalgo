package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

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
	return "Authsignal Go v1"
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

func (c Client) ValidateChallenge(request ValidateChallengeRequest) (ValidateChallengeResponse, error) {
	tokenString := request.Token
	userId := request.UserId

	claims1 := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims1, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.apiKey), nil
	})

	if err != nil {
		return ValidateChallengeResponse{}, err
	}

	if userId != claims1["UserId"].(string) {
		return ValidateChallengeResponse{}, errors.New("invalid user")
	}

	idempotencyKey := claims1["IdempotencyKey"].(string)
	action := claims1["ActionCode"].(string)

	if action != "" && idempotencyKey != "" {
		actionResult, err := c.GetAction(GetActionRequest{
			UserId:         userId,
			Action:         action,
			IdempotencyKey: idempotencyKey,
		})
		if err != nil {
			return ValidateChallengeResponse{}, err
		}

		if actionResult.State == "CHALLENGE_SUCCEEDED" {
			return ValidateChallengeResponse{true, actionResult.State}, nil
		}
	}

	return ValidateChallengeResponse{false, ""}, nil
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

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("bad request to %s, http status code of %d, status was: %s", req.URL, resp.StatusCode, resp.Status)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
