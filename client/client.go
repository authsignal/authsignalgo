package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// todo unMarshal the responses into objects.
// todo only forward body elements required along, not all requests.

const TIMEOUT = 5

type Client struct {
	apiKey      string
	apiUrl      string
	redirectUrl string
}

func New(apiKey string, apiUrl string, redirectUrl string) Client {
	c := Client{apiKey, apiUrl, redirectUrl}
	return c
}

func (c Client) defaultHeaders() http.Header {
	return http.Header{
		"Accept":       {"*/*"},
		"Content-Type": {"application/json"},
		"User-Agent":   {c.userAgent()},
	}
}

func (c Client) userAgent() string {
	return "Authsignal Python v1" // todo make module version dynamic
}

func (c Client) PrintConfig() {
	fmt.Println(c.apiUrl)
}

func (c Client) GetUser(userId string) string {
	return c.makeRequest("GET", c.apiUrl+"/v1/users/"+userId, nil)
}

func (c Client) TrackAction(request TrackRequest) TrackResponse {
	url := c.apiUrl + "/v1/users/" + request.UserId + "/actions/" + request.Action

	body, _ := json.Marshal(request)

	response := c.makeRequest("POST", url, bytes.NewBuffer(body))

	var data TrackResponse
	err := json.Unmarshal([]byte(response), &data)
	if err != nil {
		log.Fatalln(err)
		return TrackResponse{}
	}

	return data
}

func (c Client) GetAction(request GetActionRequest) string {
	url := c.apiUrl + "/v1/users/" + request.UserId + "/actions/" + request.Action + "/" + request.IdempotencyKey

	return c.makeRequest("GET", url, nil)
}

func (c Client) EnrollVerifiedAuthenticator(request EnrollVerifiedAuthenticatorRequest) string {
	url := c.apiUrl + "/v1/users/" + request.UserId + "authenticators"

	body, _ := json.Marshal(request)

	return c.makeRequest("POST", url, bytes.NewBuffer(body))
}

func (c Client) LoginWithEmail(request LoginWithEmailRequest) string {
	url := c.apiUrl + "/v1/users/email/" + request.Email + "/challenge"

	body, _ := json.Marshal(request)

	return c.makeRequest("POST", url, bytes.NewBuffer(body))
}

func (c Client) ValidateChallenge(request ValidateChallengeRequest) string {
	return "NOT YET IMPLEMENTED"
}

func (c Client) makeRequest(method string, url string, body io.Reader) string {
	client := http.Client{}
	client.Timeout = TIMEOUT
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header = c.defaultHeaders()
	req.SetBasicAuth(c.apiKey, "")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(resp.Body)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// todo return byte[] or unpack into a response object generically.
	sb := string(responseBody)
	return sb
}
