package client

type AuthsignalClient struct {
	Secret     string
	ApiBaseUrl string
}

type UserRequest struct {
	UserId string
}

type UserResponse struct {
	IsEnrolled                  bool
	Email                       string
	PhoneNumber                 string
	Username                    string
	EnrolledVerificationMethods []string
	DefaultVerificationMethod   string
}

type TrackRequest struct {
	UserId             string
	Action             string
	Email              string
	PhoneNumber        string
	IdempotencyKey     string
	RedirectUrl        string
	IpAddress          string
	UserAgent          string
	DeviceId           string
	Scope              string
	Custom             map[string]string `json:"custom"`
	RedirectToSettings bool
}

type TrackResponse struct {
	State                       string
	IdempotencyKey              string
	Url                         string
	Token                       string
	IsEnrolled                  bool
	EnrolledVerificationMethods []string
	DefaultVerificationMethod   string
}

type GetActionRequest struct {
	UserId         string
	Action         string
	IdempotencyKey string
}

type GetActionResponse struct {
	State              string
	VerificationMethod string
	CreatedAt          string
}

type EnrollVerifiedAuthenticatorRequest struct {
	UserId      string
	OobChannel  string
	PhoneNumber string
	Email       string
}

type EnrollVerifiedAuthenticatorResponse struct {
	Authenticator UserAuthenticator
	RecoveryCodes []string
}

type ValidateChallengeRequest struct {
	Token  string
	UserId string
}

type ValidateChallengeResponse struct {
	Success  bool
	State    string
	UserId   string
	Username string
}

type UserAuthenticator struct {
	UserId              string
	UserAuthenticatorId string
	AuthenticatorType   string
	CreatedAt           string
	IsDefault           bool
	VerifiedAt          string
	IsActive            bool
	OobChannel          string
	PhoneNumber         string
	Email               string
}
