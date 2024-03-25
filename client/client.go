package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

const RequestTimeout = 10 * time.Second

type Client struct {
	tenantId string
	secret   string
	apiUrl   string
	jwks     keyfunc.Keyfunc
}

func New(tenantId, secret, apiUrl string) Client {
	jwksUri := fmt.Sprintf("%s/client/public/%s/.well-known/jwks", apiUrl, tenantId)

	jwks, err := keyfunc.NewDefault([]string{jwksUri})

	if err != nil {
		log.Fatalf("Invalid tenant.\nError: %s", err)
	}

	return Client{tenantId: tenantId, secret: secret, apiUrl: apiUrl, jwks: jwks}
}

func (c Client) defaultHeaders() http.Header {
	return http.Header{
		"Accept":       {"*/*"},
		"Content-Type": {"application/json"},
		"X-Tenant-Id":  {c.tenantId},
	}
}

func (c Client) GetUser(request UserRequest) (UserResponse, error) {
	path := request.UserId
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
	parsedClaims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(request.Token, parsedClaims, func(token *jwt.Token) (interface{}, error) {
		if token.Header["alg"] == "RS256" {
			return c.jwks.Keyfunc(token)
		} else {
			return []byte(c.secret), nil
		}
	})

	if err != nil {
		return ValidateChallengeResponse{}, err
	}

	payload := parsedClaims["other"].(map[string]interface{})

	userId := payload["userId"].(string)
	username := payload["username"].(string)
	idempotencyKey := payload["idempotencyKey"].(string)
	action := payload["actionCode"].(string)

	if request.UserId != "" && request.UserId != userId {
		return ValidateChallengeResponse{}, errors.New("invalid user")
	}

	if action != "" && idempotencyKey != "" {
		actionResult, err := c.GetAction(GetActionRequest{
			UserId:         userId,
			Action:         action,
			IdempotencyKey: idempotencyKey,
		})

		if err != nil {
			return ValidateChallengeResponse{}, err
		}

		success := actionResult.State == "CHALLENGE_SUCCEEDED"

		return ValidateChallengeResponse{success, actionResult.State, userId, username}, nil
	}

	return ValidateChallengeResponse{false, "", "", ""}, nil
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
	req, err := http.NewRequest(method, fmt.Sprintf("%s/users/%s", c.apiUrl, path), body)
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
