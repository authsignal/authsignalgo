package client

// Shared
type Rule struct {
	RuleId      string      `json:"ruleId,omitempty"`
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	Metadata    interface{} `json:"metadata,omitempty"`
}

type ParsedUserAgentBrowser struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Major   string `json:"major,omitempty"`
}

type ParsedUserAgentDevice struct {
	Model  string `json:"model,omitempty"`
	Type   string `json:"type,omitempty"`
	Vendor string `json:"vendor,omitempty"`
}

type ParsedUserAgentEngine struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

type ParsedUserAgentOs struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

type ParsedUserAgentCpu struct {
	Architecture string `json:"architecture,omitempty"`
}

type ParsedUserAgent struct {
	Ua      string                  `json:"ua,omitempty"`
	Browser *ParsedUserAgentBrowser `json:"browser,omitempty"`
	Device  *ParsedUserAgentDevice  `json:"device,omitempty"`
	Engine  *ParsedUserAgentEngine  `json:"engine,omitempty"`
	Os      *ParsedUserAgentOs      `json:"os,omitempty"`
	Cpu     *ParsedUserAgentCpu     `json:"cpu,omitempty"`
}

type AaguidMapping struct {
	Name         string `json:"name,omitempty"`
	SvgIconLight string `json:"svgIconLight,omitempty"`
	SvgIconDark  string `json:"svgIconDark,omitempty"`
}

type WebAuthnCredential struct {
	CredentialId            string           `json:"credentialId,omitempty"`
	DeviceId                string           `json:"deviceId,omitempty"`
	Name                    string           `json:"name,omitempty"`
	Aaguid                  string           `json:"aaguid,omitempty"`
	AaguidMapping           *AaguidMapping   `json:"aaguidMapping,omitempty"`
	CredentialBackedUp      *bool            `json:"credentialBackedUp,omitempty"`
	CredentialDeviceType    string           `json:"credentialDeviceType,omitempty"`
	AuthenticatorAttachment string           `json:"authenticatorAttachment,omitempty"`
	ParsedUserAgent         *ParsedUserAgent `json:"parsedUserAgent,omitempty"`
}

type Crypto struct {
	Asset          string   `json:"asset,omitempty"`
	Address        string   `json:"address,omitempty"`
	TxnHash        string   `json:"txnHash,omitempty"`
	AssetAmount    *float32 `json:"assetAmount"`
	AssetAmountUsd *float32 `json:"assetAmountUsd"`
}

type TestConfig struct {
	apiSecretKey string
	apiUrl       string
}

type ActionAttributes struct {
	State string `json:"state,omitempty"`
}

// GetUser
type GetUserRequest struct {
	UserId string `json:"userId,omitempty"`
}

type GetUserResponse struct {
	IsEnrolled                  *bool       `json:"isEnrolled"`
	Email                       string      `json:"email,omitempty"`
	PhoneNumber                 string      `json:"phoneNumber,omitempty"`
	Username                    string      `json:"username,omitempty"`
	DisplayName                 string      `json:"displayName,omitempty"`
	AllowedVerificationMethods  []string    `json:"allowedVerificationMethods,omitempty"`
	EnrolledVerificationMethods []string    `json:"enrolledVerificationMethods,omitempty"`
	DefaultVerificationMethod   string      `json:"defaultVerificationMethod,omitempty"`
	Custom                      interface{} `json:"custom,omitempty"`
}

// UpdateUser
type UpdateUserRequest struct {
	UserId     string          `json:"userId,omitempty"`
	Attributes *UserAttributes `json:"attributes,omitempty"`
}

type UserAttributes struct {
	Custom      interface{} `json:"custom,omitempty"`
	Username    string      `json:"username,omitempty"`
	DisplayName string      `json:"displayName,omitempty"`
	Email       string      `json:"email,omitempty"`
	PhoneNumber string      `json:"phoneNumber,omitempty"`
}

// DeleteUser
type DeleteUserRequest struct {
	UserId string `json:"userId,omitempty"`
}

// UpdateActionState
type UpdateActionRequest struct {
	UserId         string            `json:"userId,omitempty"`
	Action         string            `json:"action,omitempty"`
	IdempotencyKey string            `json:"idempotencyKey,omitempty"`
	Attributes     *ActionAttributes `json:"attributes,omitempty"`
}

// GetAuthenticators
type GetAuthenticatorsRequest struct {
	UserId string `json:"userId,omitempty"`
}

// DeleteAuthenticator
type DeleteAuthenticatorRequest struct {
	UserId              string `json:"userId,omitempty"`
	UserAuthenticatorId string `json:"userAuthenticatorId,omitempty"`
}

// Track
type TrackRequest struct {
	UserId     string           `json:"userId,omitempty"`
	Action     string           `json:"action,omitempty"`
	Attributes *TrackAttributes `json:"attributes,omitempty"`
}

type TrackAttributes struct {
	IdempotencyKey     string      `json:"idempotencyKey,omitempty"`
	RedirectUrl        string      `json:"redirectUrl,omitempty"`
	RedirectToSettings *bool       `json:"redirectToSettings"`
	Email              string      `json:"email,omitempty"`
	PhoneNumber        string      `json:"phoneNumber,omitempty"`
	IpAddress          string      `json:"ipAddress,omitempty"`
	UserAgent          string      `json:"userAgent,omitempty"`
	Username           string      `json:"username,omitempty"`
	ChallengeId        string      `json:"challengeId,omitempty"`
	DeviceId           string      `json:"deviceId,omitempty"`
	Scope              string      `json:"scope,omitempty"`
	Custom             interface{} `json:"custom,omitempty"`
	Crypto             *Crypto     `json:"crypto,omitempty"`
}

type TrackResponse struct {
	State                       string   `json:"state,omitempty"`
	Url                         string   `json:"url,omitempty"`
	Token                       string   `json:"token,omitempty"`
	IsEnrolled                  *bool    `json:"isEnrolled"`
	IdempotencyKey              string   `json:"idempotencyKey,omitempty"`
	AllowedVerificationMethods  []string `json:"allowedVerificationMethods,omitempty"`
	EnrolledVerificationMethods []string `json:"enrolledVerificationMethods,omitempty"`
	DefaultVerificationMethod   string   `json:"defaultVerificationMethod,omitempty"`
}

// GetAction
type GetActionRequest struct {
	UserId         string `json:"userId,omitempty"`
	Action         string `json:"action,omitempty"`
	IdempotencyKey string `json:"idempotencyKey,omitempty"`
}

type GetActionResponse struct {
	State              string      `json:"state,omitempty"`
	CreatedAt          string      `json:"createdAt,omitempty"`
	StateUpdatedAt     string      `json:"stateUpdatedAt,omitempty"`
	Output             interface{} `json:"output,omitempty"`
	Rules              *[]Rule     `json:"rules,omitempty"`
	VerificationMethod string      `json:"verificationMethod,omitempty"`
}

// EnrollVerifiedAuthenticator
type EnrollVerifiedAuthenticatorRequest struct {
	UserId     string                                 `json:"userId,omitempty"`
	Attributes *EnrollVerifiedAuthenticatorAttributes `json:"attributes,omitempty"`
}

type EnrollVerifiedAuthenticatorAttributes struct {
	VerificationMethod string `json:"verificationMethod,omitempty"`
	Email              string `json:"email,omitempty"`
	PhoneNumber        string `json:"phoneNumber,omitempty"`
	OtpUri             string `json:"otpUri,omitempty"`
}

type UserAuthenticator struct {
	UserAuthenticatorId string              `json:"userAuthenticatorId,omitempty"`
	UserId              string              `json:"userId,omitempty"`
	VerificationMethod  string              `json:"verificationMethod,omitempty"`
	Email               string              `json:"email,omitempty"`
	PhoneNumber         string              `json:"phoneNumber,omitempty"`
	Username            string              `json:"username,omitempty"`
	DisplayName         string              `json:"displayName,omitempty"`
	CreatedAt           string              `json:"createdAt,omitempty"`
	VerifiedAt          string              `json:"verifiedAt,omitempty"`
	LastVerifiedAt      string              `json:"lastVerifiedAt,omitempty"`
	PreviousSmsChannel  string              `json:"previousSmsChannel,omitempty"`
	WebAuthnCredential  *WebAuthnCredential `json:"webAuthnCredential,omitempty"`
}

type EnrollVerifiedAuthenticatorResponse struct {
	Authenticator *UserAuthenticator `json:"authenticator,omitempty"`
	RecoveryCodes []string           `json:"recoveryCodes,omitempty"`
}

// ValidateChallenge
type ValidateChallengeRequest struct {
	Token  string `json:"token,omitempty"`
	UserId string `json:"userId,omitempty"`
	Action string `json:"action,omitempty"`
}

type ValidateChallengeResponse struct {
	IsValid            *bool  `json:"isValid"`
	State              string `json:"state,omitempty"`
	StateUpdatedAt     string `json:"stateUpdatedAt,omitempty"`
	UserId             string `json:"userId,omitempty"`
	Action             string `json:"action,omitempty"`
	IdempotencyKey     string `json:"idempotencyKey,omitempty"`
	VerificationMethod string `json:"verificationMethod,omitempty"`
}
