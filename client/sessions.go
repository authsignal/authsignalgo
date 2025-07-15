package client

import (
	"encoding/json"
	"strings"
)

func (c Client) CreateSession(req CreateSessionRequest) (CreateSessionResponse, error) {
	var resp CreateSessionResponse
	body, err := json.Marshal(req)
	if err != nil {
		return resp, err
	}
	data, err := c.post("/sessions", strings.NewReader(string(body)))
	if err != nil {
		return resp, err
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

func (c Client) ValidateSession(req ValidateSessionRequest) (ValidateSessionResponse, error) {
	var resp ValidateSessionResponse
	body, err := json.Marshal(req)
	if err != nil {
		return resp, err
	}
	data, err := c.post("/sessions/validate", strings.NewReader(string(body)))
	if err != nil {
		return resp, err
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

func (c Client) RefreshSession(req RefreshSessionRequest) (RefreshSessionResponse, error) {
	var resp RefreshSessionResponse
	body, err := json.Marshal(req)
	if err != nil {
		return resp, err
	}
	data, err := c.post("/sessions/refresh", strings.NewReader(string(body)))
	if err != nil {
		return resp, err
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

func (c Client) RevokeSession(req RevokeSessionRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	_, err = c.post("/sessions/revoke", strings.NewReader(string(body)))
	return err
}

func (c Client) RevokeUserSessions(req RevokeUserSessionsRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	_, err = c.post("/sessions/user/revoke", strings.NewReader(string(body)))
	return err
}
