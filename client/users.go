package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
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

func (c Client) QueryUsers(input QueryUsersRequest) (QueryUsersResponse, error) {
	params := url.Values{}

	if input.Username != "" {
		params.Set("username", input.Username)
	}
	if input.Email != "" {
		params.Set("email", input.Email)
	}
	if input.PhoneNumber != "" {
		params.Set("phoneNumber", input.PhoneNumber)
	}
	if input.Token != "" {
		params.Set("token", input.Token)
	}
	if input.Limit != nil {
		params.Set("limit", strconv.Itoa(*input.Limit))
	}
	if input.LastEvaluatedUserId != "" {
		params.Set("lastEvaluatedUserId", input.LastEvaluatedUserId)
	}

	path := "/users"
	if len(params) > 0 {
		path = fmt.Sprintf("/users?%s", params.Encode())
	}

	response, err := c.get(path)
	if err != nil {
		return QueryUsersResponse{}, err
	}

	var data QueryUsersResponse
	err = json.Unmarshal(response, &data)
	if err != nil {
		return QueryUsersResponse{}, err
	}

	return data, nil
}
