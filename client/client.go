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
// todo only forward body elements required along, not all of request.

type Client struct {
	Api_key      string
	Api_url      string
	Redirect_url string
}

func New(api_key string, api_url string, redirect_url string) Client {
	c := Client{api_key, api_url, redirect_url}
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
	fmt.Println(c.Api_url)
}

func (c Client) GetUser(user_id string) string {
	return c.makeRequest("GET", c.Api_url+"/v1/users/"+user_id, nil)
}

func (c Client) TrackAction(request TrackRequest) string {
	url := c.Api_url + "/v1/users/" + request.UserId + "/actions/" + request.Action

	body, _ := json.Marshal(request)

	return c.makeRequest("POST", url, bytes.NewBuffer(body))
}

func (c Client) GetAction(request GetActionRequest) string {
	url := c.Api_url + "/v1/users/" + request.UserId + "/actions/" + request.Action + "/" + request.IdempotencyKey

	return c.makeRequest("GET", url, nil)
}

func (c Client) EnrollVerifiedAuthenticator(request EnrollVerifiedAuthenticatorRequest) string {
	url := c.Api_url + "/v1/users/" + request.UserId + "authenticators"

	body, _ := json.Marshal(request)

	return c.makeRequest("POST", url, bytes.NewBuffer(body))
}

func (c Client) LoginWithEmail(request LoginWithEmailRequest) string {
	url := c.Api_url + "/v1/users/email/" + request.Email + "/challenge"

	body, _ := json.Marshal(request)

	return c.makeRequest("POST", url, bytes.NewBuffer(body))
}

func (c Client) ValidateChallenge(request ValidateChallengeRequest) string {
	return "NOT YET IMPLEMENTED"
}

func (c Client) makeRequest(method string, url string, body io.Reader) string {
	client := http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header = c.defaultHeaders()
	req.SetBasicAuth(c.Api_key, "")

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

	sb := string(responseBody)
	return sb
}
