package client

import (
	"bytes"
	"encoding/json"
)

func (c Client) StartFlow(input StartFlowRequest) (StartFlowResponse, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return StartFlowResponse{}, err
	}
	response, err := c.post("/flows", bytes.NewBuffer(body))
	if err != nil {
		return StartFlowResponse{}, err
	}
	var data StartFlowResponse
	if err = json.Unmarshal(response, &data); err != nil {
		return StartFlowResponse{}, err
	}
	return data, nil
}

func (c Client) VerifyFlow(input VerifyFlowRequest) (VerifyFlowResponse, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return VerifyFlowResponse{}, err
	}
	response, err := c.post("/flows/verify", bytes.NewBuffer(body))
	if err != nil {
		return VerifyFlowResponse{}, err
	}
	var data VerifyFlowResponse
	if err = json.Unmarshal(response, &data); err != nil {
		return VerifyFlowResponse{}, err
	}
	return data, nil
}
