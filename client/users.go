package client

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (c Client) GetUser(input GetUserRequest) (GetUserResponse, error) {
	path := fmt.Sprintf("/users/%s", input.UserId)
	response, err := c.get(path)
	if err != nil {
		return GetUserResponse{}, err
	}

	var data GetUserResponse
	err = json.Unmarshal(response, &data)
	if err != nil {
		return GetUserResponse{}, err
	}

	return data, nil
}

func (c Client) UpdateUser(input UpdateUserRequest) (UserAttributes, error) {
	body, err := json.Marshal(input.Attributes)
	if err != nil {
		return UserAttributes{}, err
	}

	path := fmt.Sprintf("/users/%s", input.UserId)
	response, err := c.patch(path, bytes.NewBuffer(body))
	if err != nil {
		return UserAttributes{}, err
	}

	var data UserAttributes
	err = json.Unmarshal(response, &data)
	if err != nil {
		return UserAttributes{}, err
	}

	return data, nil
}

func (c Client) DeleteUser(input DeleteUserRequest) error {
	path := fmt.Sprintf("/users/%s", input.UserId)
	_, err := c.delete(path)
	if err != nil {
		return err
	}

	return nil
}
