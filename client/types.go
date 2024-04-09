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
	Email              string            `json:"email,omitempty"`
	PhoneNumber        string            `json:"phoneNumber,omitempty"`
	IdempotencyKey     string            `json:"idempotencyKey,omitempty"`
	RedirectUrl        string            `json:"redirectUrl,omitempty"`
	IpAddress          string            `json:"ipAddress,omitempty"`
	UserAgent          string            `json:"userAgent,omitempty"`
	DeviceId           string            `json:"deviceId,omitempty"`
	Scope              string            `json:"scope,omitempty"`
	Custom             map[string]string `json:"custom,omitempty"`
	RedirectToSettings bool              `json:"redirectToSettings,omitempty"`
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
	OobChannel  string `json:"oobChannel"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	Email       string `json:"email,omitempty"`
}

type EnrollVerifiedAuthenticatorResponse struct {
	Authenticator UserAuthenticator
	RecoveryCodes []string
}

type ValidateChallengeRequest struct {
	Token  string `json:"token"`
	UserId string `json:"userId"`
}

type ValidateChallengeResponse struct {
	IsValid            bool
	State              string
	StateUpdatedAt     string
	UserId             string
	Action             string `json:"actionCode"`
	IdempotencyKey     string
	VerificationMethod string
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

type AuthsignalApiError struct {
	Error            string
	ErrorDescription string
}
