package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func (c Client) Challenge(req ChallengeRequest) (ChallengeResponse, error) {
	var resp ChallengeResponse
	body, err := json.Marshal(req)
	if err != nil {
		return resp, err
	}
	data, err := c.post("/challenge", strings.NewReader(string(body)))
	if err != nil {
		return resp, err
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

func (c Client) Verify(req VerifyRequest) (VerifyResponse, error) {
	var resp VerifyResponse
	body, err := json.Marshal(req)
	if err != nil {
		return resp, err
	}
	data, err := c.post("/verify", strings.NewReader(string(body)))
	if err != nil {
		return resp, err
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

func (c Client) ClaimChallenge(req ClaimChallengeRequest) (ClaimChallengeResponse, error) {
	var resp ClaimChallengeResponse
	body, err := json.Marshal(req)
	if err != nil {
		return resp, err
	}
	data, err := c.post("/claim", strings.NewReader(string(body)))
	if err != nil {
		return resp, err
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

func (c Client) GetChallenge(req GetChallengeRequest) (GetChallengeResponse, error) {
	var resp GetChallengeResponse
	// Build query string
	u, err := url.Parse(fmt.Sprintf("%s/challenges", c.ApiUrl))
	if err != nil {
		return resp, err
	}
	q := u.Query()
	if req.ChallengeId != "" {
		q.Set("challengeId", req.ChallengeId)
	}
	if req.UserId != "" {
		q.Set("userId", req.UserId)
	}
	if req.Action != "" {
		q.Set("action", req.Action)
	}
	if req.VerificationMethod != "" {
		q.Set("verificationMethod", string(req.VerificationMethod))
	}
	u.RawQuery = q.Encode()

	data, err := c.get(u.Path + "?" + u.RawQuery)
	if err != nil {
		return resp, err
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}
