package client

type StartFlowRequest struct {
	ActionCode  string               `json:"actionCode"`
	User        *UserLookup          `json:"user,omitempty"`
	Attributes  *ChallengeAttributes `json:"attributes,omitempty"`
	RedirectUrl *string              `json:"redirectUrl,omitempty"`
	ClientId    *string              `json:"clientId,omitempty"`
}

type StartFlowResponse struct {
	Action         FlowAction `json:"action"`
	ChallengeToken string     `json:"challengeToken"`
	ChallengeUrl   string     `json:"challengeUrl"`
	User           *FlowUser  `json:"user,omitempty"`
}

type VerifyFlowRequest struct {
	ActionCode     string `json:"actionCode"`
	ChallengeToken string `json:"challengeToken"`
}

type VerifyFlowResponse struct {
	Action  FlowAction             `json:"action"`
	Session *AuthenticationSession `json:"session,omitempty"`
	User    *FlowUser              `json:"user,omitempty"`
}

type UserLookup struct {
	UserId      *string `json:"userId,omitempty"`
	Email       *string `json:"email,omitempty"`
	PhoneNumber *string `json:"phoneNumber,omitempty"`
	Username    *string `json:"username,omitempty"`
}

type FlowAction struct {
	State          FlowState             `json:"state"`
	CompletedSteps []CompletedActionStep `json:"completedSteps"`
	NextStep       *ActionStep           `json:"nextStep,omitempty"`
}

type ActionStep struct {
	StepType            ActionStepType `json:"stepType"`
	VerificationMethods []string       `json:"verificationMethods"`
}

type CompletedActionStep struct {
	StepType            ActionStepType `json:"stepType"`
	VerificationMethod  string         `json:"verificationMethod"`
	UserAuthenticatorId *string        `json:"userAuthenticatorId,omitempty"`
}

type ChallengeAttributes struct {
	DeviceId  *string                `json:"deviceId,omitempty"`
	IpAddress *string                `json:"ipAddress,omitempty"`
	UserAgent *string                `json:"userAgent,omitempty"`
	Custom    map[string]interface{} `json:"custom,omitempty"`
	Locale    *string                `json:"locale,omitempty"`
}

type FlowUser struct {
	UserId         string                  `json:"userId"`
	Authenticators []FlowUserAuthenticator `json:"authenticators"`
	Username       *string                 `json:"username,omitempty"`
}

type FlowUserAuthenticator struct {
	UserAuthenticatorId string  `json:"userAuthenticatorId"`
	VerificationMethod  string  `json:"verificationMethod"`
	Email               *string `json:"email,omitempty"`
	PhoneNumber         *string `json:"phoneNumber,omitempty"`
	Username            *string `json:"username,omitempty"`
	DisplayName         *string `json:"displayName,omitempty"`
}

type AuthenticationSession struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type FlowState string

const (
	FlowStateChallengeRequired  FlowState = "CHALLENGE_REQUIRED"
	FlowStateChallengeSucceeded FlowState = "CHALLENGE_SUCCEEDED"
	FlowStateChallengeFailed    FlowState = "CHALLENGE_FAILED"
)

type ActionStepType string

const (
	ActionStepTypeVerificationRequired ActionStepType = "VERIFICATION_REQUIRED"
	ActionStepTypeEnrollmentRequired   ActionStepType = "ENROLLMENT_REQUIRED"
	ActionStepTypeEnrollmentOptional   ActionStepType = "ENROLLMENT_OPTIONAL"
)
