package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Client struct {
	Api_key string
	Api_url string
	Timeout int    //todo set default, better datatype.
	Version int    //todo set default, better datatype.
	Session string //todo find why this is passed.
}

func New(api_key string, api_url string, timeout int, version int, session string) Client {
	c := Client{api_key, api_url, timeout, version, session}
	return c
}

func (c Client) PrintConfig() {
	fmt.Println(c.Api_url)
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

func (c Client) GetUser(user_id string) string {
	client := http.Client{}
	req, err := http.NewRequest("GET", c.Api_url+"/v1/users/"+user_id, nil)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)
	return sb
}
