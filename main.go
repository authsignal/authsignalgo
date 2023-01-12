package main

import (
	"authsignalgo/client"
	"fmt"
)

func main() {
	c := client.Client{Api_url: "https://dev-signal.authsignal.com", Api_key: "mURA8eRlODYiEd+butwzGp/5GTaC3tyEoFcufvVxtd0OoJYVK9EVuA==", Version: 1, Timeout: 2, Session: ""}
	c.PrintConfig()
	fmt.Println(c.GetUser("1"))
}
