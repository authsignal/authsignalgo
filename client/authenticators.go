package client

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (c Client) EnrollVerifiedAuthenticator(input EnrollVerifiedAuthenticatorRequest) (EnrollVerifiedAuthenticatorResponse, error) {
	body, err := json.Marshal(input.Attributes)
	if err != nil {
		return EnrollVerifiedAuthenticatorResponse{}, err
	}

	path := fmt.Sprintf("/users/%s/authenticators", input.UserId)
	response, err := c.post(path, bytes.NewBuffer(body))
	if err != nil {
		return EnrollVerifiedAuthenticatorResponse{}, err
	}

	var data EnrollVerifiedAuthenticatorResponse
	err = json.Unmarshal(response, &data)
	if err != nil {
		return EnrollVerifiedAuthenticatorResponse{}, err
	}

	return data, nil
}

func (c Client) GetAuthenticators(input GetAuthenticatorsRequest) ([]UserAuthenticator, error) {
	path := fmt.Sprintf("/users/%s/authenticators", input.UserId)
	response, err := c.get(path)
	if err != nil {
		return []UserAuthenticator{}, err
	}

	var data []UserAuthenticator
	err = json.Unmarshal(response, &data)
	if err != nil {
		return []UserAuthenticator{}, err
	}

	return data, nil
}

func (c Client) DeleteAuthenticator(input DeleteAuthenticatorRequest) error {
	path := fmt.Sprintf("/users/%s/authenticators/%s", input.UserId, input.UserAuthenticatorId)
	_, err := c.delete(path)
	if err != nil {
		return err
	}

	return nil
}
