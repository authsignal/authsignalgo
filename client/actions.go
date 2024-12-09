package client

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (c Client) Track(input TrackRequest) (TrackResponse, error) {
	var body []byte
	var err error
	if input.Attributes == nil {
		body = []byte("{}")
	} else {
		body, err = json.Marshal(input.Attributes)
		if err != nil {
			return TrackResponse{}, err
		}
	}

	path := fmt.Sprintf("/users/%s/actions/%s", input.UserId, input.Action)
	response, err := c.post(path, bytes.NewBuffer(body))
	if err != nil {
		return TrackResponse{}, err
	}

	var data TrackResponse
	err = json.Unmarshal(response, &data)
	if err != nil {
		return TrackResponse{}, err
	}

	return data, nil
}

func (c Client) GetAction(input GetActionRequest) (GetActionResponse, error) {
	path := fmt.Sprintf("/users/%s/actions/%s/%s", input.UserId, input.Action, input.IdempotencyKey)
	response, err := c.get(path)
	if err != nil {
		return GetActionResponse{}, err
	}

	var data GetActionResponse
	err = json.Unmarshal(response, &data)
	if err != nil {
		return GetActionResponse{}, err
	}

	return data, nil
}

func (c Client) UpdateAction(input UpdateActionRequest) (ActionAttributes, error) {
	body, err := json.Marshal(input.Attributes)
	if err != nil {
		return ActionAttributes{}, err
	}

	path := fmt.Sprintf("/users/%s/actions/%s/%s", input.UserId, input.Action, input.IdempotencyKey)
	response, err := c.patch(path, bytes.NewBuffer(body))
	if err != nil {
		return ActionAttributes{}, err
	}

	var data ActionAttributes
	err = json.Unmarshal(response, &data)
	if err != nil {
		return ActionAttributes{}, err
	}

	return data, nil
}

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
