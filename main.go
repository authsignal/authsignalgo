package main

import (
	"authsignalgo/client"
	"fmt"
)

func main() {
	c := client.Client{
		Api_url:      "https://dev-signal.authsignal.com",
		Api_key:      "mURA8eRlODYiEd+butwzGp/5GTaC3tyEoFcufvVxtd0OoJYVK9EVuA==",
		Redirect_url: "",
	}
	c.PrintConfig()
	fmt.Println(c.GetUser("1"))

	fmt.Println()
	fmt.Println("Enroll Verified Authenticator")
	enrollVerifiedAuthenticatorRequest := client.EnrollVerifiedAuthenticatorRequest{
		UserId:      "1",
		OobChannel:  "blah",
		PhoneNumber: "024525252",
		Email:       "test@email.me",
	}
	fmt.Println(c.EnrollVerifiedAuthenticator(enrollVerifiedAuthenticatorRequest))

	fmt.Println()
	fmt.Println("Login With Email")
	loginWithEmailRequest := client.LoginWithEmailRequest{
		Email:       "test@email.me",
		RedirectUrl: "http://www.authsignal.com",
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
