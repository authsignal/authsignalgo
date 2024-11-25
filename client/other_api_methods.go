package client

import (
	"bytes"
	"encoding/json"
)

func (c Client) ValidateChallenge(input ValidateChallengeRequest) (ValidateChallengeResponse, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return ValidateChallengeResponse{}, err
	}

	path := "/validate"
	response, err := c.post(path, bytes.NewBuffer(body))
	if err != nil {
		return ValidateChallengeResponse{}, err
	}

	var data ValidateChallengeResponse
	err = json.Unmarshal(response, &data)
	if err != nil {
		return ValidateChallengeResponse{}, err
	}

	return data, nil
}
