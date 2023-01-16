package main

import (
	"authsignalgo/client"
	"fmt"
)

func main() {
	c := client.New(
		"https://dev-signal.authsignal.com",
		"mURA8eRlODYiEd+butwzGp/5GTaC3tyEoFcufvVxtd0OoJYVK9EVuA==",
		"")

	fmt.Println()
	fmt.Println("Get User")
	fmt.Println(c.GetUser("1"))

	fmt.Println()
	fmt.Println("Enroll Verified Authenticator")
	enrollVerifiedAuthenticatorRequest := client.EnrollVerifiedAuthenticatorRequest{
		UserId:      "1",
		OobChannel:  "blah",
		PhoneNumber: "024525252",
		Email:       "test@email.me",
	}
	response, err := c.EnrollVerifiedAuthenticator(enrollVerifiedAuthenticatorRequest)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)

	fmt.Println()
	fmt.Println("Login With Email")
	loginWithEmailRequest := client.LoginWithEmailRequest{
		Email:       "test@email.me",
		RedirectUrl: "https://www.authsignal.com",
	}
	fmt.Println(c.LoginWithEmail(loginWithEmailRequest))

	fmt.Println()
	fmt.Println("Tracking Request")
	trackRequest := client.TrackRequest{
		UserId: "1",
		Action: "ABC",
	}
	fmt.Println(c.TrackAction(trackRequest))

	fmt.Println()
	fmt.Println("Get Action Request")
	getActionRequest := client.GetActionRequest{
		UserId:         "1",
		Action:         "ABC",
		IdempotencyKey: "ACBD12312",
	}
	fmt.Println(c.GetAction(getActionRequest))
}
