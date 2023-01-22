package main

import (
	"authsignalgo/client"
	"fmt"
	"github.com/golang-jwt/jwt"
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
		OobChannel:  "EMAIL_MAGIC_LINK",
		PhoneNumber: "024525252",
		Email:       "test@email.me",
	}
	fmt.Println(c.EnrollVerifiedAuthenticator(enrollVerifiedAuthenticatorRequest))

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
	action, _ := c.TrackAction(trackRequest)

	key := action.IdempotencyKey
	fmt.Println(action)

	fmt.Println()
	fmt.Println("Get Action Request")
	getActionRequest := client.GetActionRequest{
		UserId:         "1",
		Action:         "ABC",
		IdempotencyKey: key,
	}
	fmt.Println(c.GetAction(getActionRequest))

	fmt.Println()
	fmt.Println("Validate Challenge")
	tokenString, _ := generateJWT("1", action.IdempotencyKey)
	validateChallengeRequest := client.ValidateChallengeRequest{
		UserId: "1",
		Token:  tokenString,
	}
	fmt.Println(c.ValidateChallenge(validateChallengeRequest))

}

func generateJWT(userId, idempotencyKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"UserId":         userId,
		"TenantId":       12342,
		"PublishableKey": "1",
		"Email":          "email@test.me",
		"PhoneNumber":    "242525252",
		"ActionCode":     "ABC",
		"IdempotencyKey": idempotencyKey,
	})

	tokenString, err := token.SignedString([]byte("mURA8eRlODYiEd+butwzGp/5GTaC3tyEoFcufvVxtd0OoJYVK9EVuA=="))
	if err != nil {
		return "Signing Error", err
	}
	return tokenString, nil
}
