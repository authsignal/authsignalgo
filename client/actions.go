package client

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (c Client) Track(input TrackRequest) (TrackResponse, error) {
	body, err := json.Marshal(input.Attributes)
	if err != nil {
		return TrackResponse{}, err
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

func (c Client) UpdateAction(input UpdateActionRequest) (UpdateActionResponse, error) {
	body, err := json.Marshal(input.Attributes)
	if err != nil {
		return UpdateActionResponse{}, err
	}

	path := fmt.Sprintf("/users/%s/actions/%s/%s", input.UserId, input.Action, input.IdempotencyKey)
	response, err := c.patch(path, bytes.NewBuffer(body))
	if err != nil {
		return UpdateActionResponse{}, err
	}

	var data UpdateActionResponse
	err = json.Unmarshal(response, &data)
	if err != nil {
		return UpdateActionResponse{}, err
	}

	return data, nil
}
