package client

import "fmt"

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
