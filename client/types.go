package client

type AuthsignalConstructor struct {
	Secret      string
	ApiBaseUrl  string //TODO APIBaseURL
	RedirectUrl string
}

type UserRequest struct {
	UserId string // todo ID
}

type UserResponse struct {
	IsEnrolled bool
}

type TrackRequest struct {
	UserId             string
	Action             string
	Email              string
	IdempotencyKey     string
	RedirectUrl        string
	IpAddress          string
	UserAgent          string
	DeviceId           string
	RedirectToSettings bool
}

type TrackResponse struct {
	State          string
	IdempotencyKey string
	RuleIds        []string
	Url            string
	IsEnrolled     bool
	ChallengeUrl   string
}

type GetActionRequest struct {
	UserId         string
	Action         string
	IdempotencyKey string
}

type GetActionResponse struct {
	State string
}

type EnrollVerifiedAuthenticatorRequest struct {
	UserId      string `json:"userId"`
	OobChannel  string `json:"oobChannel"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
}

type EnrollVerifiedAuthenticatorResponse struct {
	Authenticator UserAuthenticator
	RecoveryCodes []string
}

type LoginWithEmailRequest struct {
	Email       string
	RedirectUrl string
}

type LoginWithEmailResponse struct {
	Url string
}

type ValidateChallengeRequest struct {
	UserId string
	Token  string
}

type ValidateChallengeResponse struct {
	Success bool
	State   string
}

type UserAuthenticator struct {
	UserId              string
	UserAuthenticatorId string
	AuthenticatorType   AuthenticatorType
	CreatedAt           string
	IsDefault           bool
	VerifiedAt          string
	IsActive            bool
	OobChannel          OobChannel
	OtpBinding          OtpBinding
	PhoneNumber         string
	Email               string
}

type AuthenticatorType int

const (
	OOB = iota
	OTP
)

type OtpBinding struct {
	Secret string
	Uri    string
}
type OobChannel int

const (
	SMS = iota
	EMAIL_MAGIC_LINK
)

type RedirectTokenPayload struct {
	TenantId       string
	PublishableKey string
	UserId         string
	Email          string
	PhoneNumber    string
	ActionCode     string
	IdempotencyKey string
}
