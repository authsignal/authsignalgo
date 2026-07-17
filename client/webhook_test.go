package client

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
	"time"
)

const testSecretKey = "test-secret-key-123"

func getTestWebhook() *Webhook {
	return NewWebhook(testSecretKey)
}

func generateTestSignature(payload string, timestamp int64, secret string) string {
	hmacContent := fmt.Sprintf("%d.%s", timestamp, payload)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(hmacContent))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	signature = strings.ReplaceAll(signature, "=", "")
	return fmt.Sprintf("t=%d,v2=%s", timestamp, signature)
}

func TestInvalidSignatureFormat(t *testing.T) {
	webhook := getTestWebhook()
	payload := "{}"
	signature := "123"

	_, err := webhook.ConstructEvent(payload, signature, DefaultTolerance)

	if err == nil {
		t.Error("Expected an error to be returned")
		return
	}

	invalidSigErr, ok := err.(*InvalidSignatureError)
	if !ok {
		t.Errorf("Expected InvalidSignatureError, got %T", err)
		return
	}

	if invalidSigErr.Message != "Signature format is invalid." {
		t.Errorf("Expected 'Signature format is invalid.', got '%s'", invalidSigErr.Message)
	}
}

func TestTimestampToleranceError(t *testing.T) {
	webhook := getTestWebhook()
	payload := "{}"
	signature := "t=1630000000,v2=invalid_signature"

	_, err := webhook.ConstructEvent(payload, signature, DefaultTolerance)

	if err == nil {
		t.Error("Expected an error to be returned")
		return
	}

	invalidSigErr, ok := err.(*InvalidSignatureError)
	if !ok {
		t.Errorf("Expected InvalidSignatureError, got %T", err)
		return
	}

	if invalidSigErr.Message != "Timestamp is outside the tolerance zone." {
		t.Errorf("Expected 'Timestamp is outside the tolerance zone.', got '%s'", invalidSigErr.Message)
	}
}

func TestInvalidComputedSignature(t *testing.T) {
	webhook := getTestWebhook()
	payload := "{}"
	timestamp := time.Now().Unix()
	signature := fmt.Sprintf("t=%d,v2=invalid_signature", timestamp)

	_, err := webhook.ConstructEvent(payload, signature, DefaultTolerance)

	if err == nil {
		t.Error("Expected an error to be returned")
		return
	}

	invalidSigErr, ok := err.(*InvalidSignatureError)
	if !ok {
		t.Errorf("Expected InvalidSignatureError, got %T", err)
		return
	}

	if invalidSigErr.Message != "Signature mismatch." {
		t.Errorf("Expected 'Signature mismatch.', got '%s'", invalidSigErr.Message)
	}
}

func TestValidSignature(t *testing.T) {
	webhook := getTestWebhook()

	payload := `{"version":1,"id":"bc1598bc-e5d6-4c69-9afb-1a6fe3469d6e","source":"https://authsignal.com","time":"2025-02-20T01:51:56.070Z","tenantId":"7752d28e-e627-4b1b-bb81-b45d68d617bc","type":"email.created","data":{"to":"chris@authsignal.com","code":"157743","userId":"b9f74d36-fcfc-4efc-87f1-3664ab5a7fb0","actionCode":"accountRecovery","idempotencyKey":"ba8c1a7c-775d-4dff-9abe-be798b7b8bb9","verificationMethod":"EMAIL_OTP"}}`

	tolerance := -1

	timestamp := time.Now().Unix()
	signature := generateTestSignature(payload, timestamp, testSecretKey)

	event, err := webhook.ConstructEvent(payload, signature, tolerance)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if event == nil {
		t.Error("Expected event to be returned")
		return
	}

	if event.Version != 1 {
		t.Errorf("Expected version 1, got %d", event.Version)
	}

	actionCode, ok := event.Data["actionCode"].(string)
	if !ok {
		t.Error("Expected actionCode to be a string")
		return
	}

	if actionCode != "accountRecovery" {
		t.Errorf("Expected actionCode 'accountRecovery', got '%s'", actionCode)
	}
}

func TestValidSignatureWhenTwoApiKeysActive(t *testing.T) {
	webhook := getTestWebhook()

	payload := `{"version":1,"id":"af7be03c-ea8f-4739-b18e-8b48fcbe4e38","source":"https://authsignal.com","time":"2025-02-20T01:47:17.248Z","tenantId":"7752d28e-e627-4b1b-bb81-b45d68d617bc","type":"email.created","data":{"to":"chris@authsignal.com","code":"718190","userId":"b9f74d36-fcfc-4efc-87f1-3664ab5a7fb0","actionCode":"accountRecovery","idempotencyKey":"68d68190-fac9-4e91-b277-c63d31d3c6b1","verificationMethod":"EMAIL_OTP"}}`

	tolerance := -1

	timestamp := time.Now().Unix()
	validSignature := generateTestSignature(payload, timestamp, testSecretKey)
	signature := validSignature + ",v2=oldKeyInvalidSignature123"

	event, err := webhook.ConstructEvent(payload, signature, tolerance)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if event == nil {
		t.Error("Expected event to be returned")
	}
}

func TestEventWithCustomVariables(t *testing.T) {
	webhook := getTestWebhook()
	payload := `{"version":1,"id":"bc1598bc-e5d6-4c69-9afb-1a6fe3469d6e","source":"https://authsignal.com","time":"2025-02-20T01:51:56.070Z","tenantId":"7752d28e-e627-4b1b-bb81-b45d68d617bc","type":"sms.created","data":{"actionCode":"smsVerify","customVariables":{"action_journeyType":"ForgotChangePassword","retryCount":2,"isRecovery":true,"channels":["sms","email"]}}}`
	timestamp := time.Now().Unix()

	event, err := webhook.ConstructEvent(payload, generateTestSignature(payload, timestamp, testSecretKey), DefaultTolerance)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	customVariables, ok := event.Data["customVariables"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected customVariables to be an object, got %T", event.Data["customVariables"])
	}

	if customVariables["action_journeyType"] != "ForgotChangePassword" {
		t.Errorf("Expected journey type to be preserved, got %v", customVariables["action_journeyType"])
	}
	if customVariables["retryCount"] != float64(2) {
		t.Errorf("Expected numeric value to be preserved, got %v", customVariables["retryCount"])
	}
	if customVariables["isRecovery"] != true {
		t.Errorf("Expected boolean value to be preserved, got %v", customVariables["isRecovery"])
	}
	channels, ok := customVariables["channels"].([]interface{})
	if !ok || len(channels) != 2 || channels[0] != "sms" || channels[1] != "email" {
		t.Errorf("Expected array value to be preserved, got %v", customVariables["channels"])
	}
}

func TestLogEventBatch(t *testing.T) {
	webhook := getTestWebhook()
	payload := `{"records":[{"version":1,"id":"bc1598bc-e5d6-4c69-9afb-1a6fe3469d6e","source":"https://authsignal.com","time":"2025-02-20T01:51:56.070Z","tenantId":"7752d28e-e627-4b1b-bb81-b45d68d617bc","type":"action.log_created","record":{"userId":"b9f74d36-fcfc-4efc-87f1-3664ab5a7fb0","customVariables":{"journeyType":"accountRecovery"}}}]}`
	timestamp := time.Now().Unix()

	batch, err := webhook.ConstructLogEventBatchWithDefaultTolerance(
		payload,
		generateTestSignature(payload, timestamp, testSecretKey),
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(batch.Records) != 1 {
		t.Fatalf("Expected one record, got %d", len(batch.Records))
	}
	customVariables, ok := batch.Records[0].Record["customVariables"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected customVariables to be an object, got %T", batch.Records[0].Record["customVariables"])
	}
	if customVariables["journeyType"] != "accountRecovery" {
		t.Errorf("Expected journey type to be preserved, got %v", customVariables["journeyType"])
	}
}

func TestLogEventBatchPassedToConstructEvent(t *testing.T) {
	webhook := getTestWebhook()
	payload := `{"records":[]}`
	timestamp := time.Now().Unix()

	_, err := webhook.ConstructEvent(
		payload,
		generateTestSignature(payload, timestamp, testSecretKey),
		DefaultTolerance,
	)
	if err == nil {
		t.Fatal("Expected an error to be returned")
	}
	if err.Error() != "Payload is a batch of log events. Use ConstructLogEventBatch instead." {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestValidSignatureWithOldKeyFirst(t *testing.T) {
	webhook := getTestWebhook()

	payload := `{"version":1,"id":"test-id","source":"https://authsignal.com","time":"2025-02-20T01:47:17.248Z","tenantId":"test-tenant","type":"email.created","data":{}}`

	tolerance := -1

	timestamp := time.Now().Unix()
	hmacContent := fmt.Sprintf("%d.%s", timestamp, payload)
	mac := hmac.New(sha256.New, []byte(testSecretKey))
	mac.Write([]byte(hmacContent))
	validSig := strings.ReplaceAll(base64.StdEncoding.EncodeToString(mac.Sum(nil)), "=", "")

	signature := fmt.Sprintf("t=%d,v2=invalidOldKeySignature,v2=%s", timestamp, validSig)

	event, err := webhook.ConstructEvent(payload, signature, tolerance)

	if err != nil {
		t.Errorf("Expected no error when valid signature is second, got %v", err)
		return
	}

	if event == nil {
		t.Error("Expected event to be returned")
	}
}

func TestEmptySignature(t *testing.T) {
	webhook := getTestWebhook()
	payload := "{}"
	signature := ""

	_, err := webhook.ConstructEvent(payload, signature, DefaultTolerance)

	if err == nil {
		t.Error("Expected an error to be returned")
		return
	}

	invalidSigErr, ok := err.(*InvalidSignatureError)
	if !ok {
		t.Errorf("Expected InvalidSignatureError, got %T", err)
		return
	}

	if invalidSigErr.Message != "Signature format is invalid." {
		t.Errorf("Expected 'Signature format is invalid.', got '%s'", invalidSigErr.Message)
	}
}

func TestConstructEventWithDefaultTolerance(t *testing.T) {
	webhook := getTestWebhook()
	payload := "{}"
	signature := "t=1630000000,v2=invalid_signature"

	_, err := webhook.ConstructEventWithDefaultTolerance(payload, signature)

	if err == nil {
		t.Error("Expected an error to be returned")
		return
	}

	invalidSigErr, ok := err.(*InvalidSignatureError)
	if !ok {
		t.Errorf("Expected InvalidSignatureError, got %T", err)
		return
	}

	if invalidSigErr.Message != "Timestamp is outside the tolerance zone." {
		t.Errorf("Expected 'Timestamp is outside the tolerance zone.', got '%s'", invalidSigErr.Message)
	}
}

func TestConstructEventWithDefaultToleranceValid(t *testing.T) {
	webhook := getTestWebhook()
	payload := `{"version":1,"type":"test.event","id":"123","source":"test","time":"2025-01-01T00:00:00Z","tenantId":"tenant","data":{}}`

	timestamp := time.Now().Unix()
	signature := generateTestSignature(payload, timestamp, testSecretKey)

	event, err := webhook.ConstructEventWithDefaultTolerance(payload, signature)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if event == nil {
		t.Error("Expected event to be returned")
		return
	}

	if event.Type != "test.event" {
		t.Errorf("Expected type 'test.event', got '%s'", event.Type)
	}
}

func TestMissingTimestamp(t *testing.T) {
	webhook := getTestWebhook()
	payload := "{}"
	signature := "v2=someSignature"

	_, err := webhook.ConstructEvent(payload, signature, DefaultTolerance)

	if err == nil {
		t.Error("Expected an error to be returned")
		return
	}

	invalidSigErr, ok := err.(*InvalidSignatureError)
	if !ok {
		t.Errorf("Expected InvalidSignatureError, got %T", err)
		return
	}

	if invalidSigErr.Message != "Signature format is invalid." {
		t.Errorf("Expected 'Signature format is invalid.', got '%s'", invalidSigErr.Message)
	}
}

func TestMissingSignature(t *testing.T) {
	webhook := getTestWebhook()
	payload := "{}"
	signature := "t=1234567890"

	_, err := webhook.ConstructEvent(payload, signature, DefaultTolerance)

	if err == nil {
		t.Error("Expected an error to be returned")
		return
	}

	invalidSigErr, ok := err.(*InvalidSignatureError)
	if !ok {
		t.Errorf("Expected InvalidSignatureError, got %T", err)
		return
	}

	if invalidSigErr.Message != "Signature format is invalid." {
		t.Errorf("Expected 'Signature format is invalid.', got '%s'", invalidSigErr.Message)
	}
}

func TestInvalidJSON(t *testing.T) {
	webhook := getTestWebhook()
	payload := "not valid json"

	timestamp := time.Now().Unix()
	signature := generateTestSignature(payload, timestamp, testSecretKey)

	_, err := webhook.ConstructEvent(payload, signature, -1)

	if err == nil {
		t.Error("Expected an error for invalid JSON")
		return
	}

	_, ok := err.(*InvalidPayloadError)
	if !ok {
		t.Errorf("Expected InvalidPayloadError, got %T", err)
	}
}

func TestTimestampAtExactTolerance(t *testing.T) {
	webhook := getTestWebhook()
	payload := `{"version":1,"type":"test","id":"1","source":"test","time":"2025-01-01T00:00:00Z","tenantId":"t","data":{}}`

	timestamp := time.Now().Unix() - (DefaultTolerance * 60)
	signature := generateTestSignature(payload, timestamp, testSecretKey)

	event, err := webhook.ConstructEvent(payload, signature, DefaultTolerance)

	if err != nil {
		t.Errorf("Expected no error at exact tolerance boundary, got %v", err)
		return
	}

	if event == nil {
		t.Error("Expected event to be returned")
	}
}

func TestTimestampJustOutsideTolerance(t *testing.T) {
	webhook := getTestWebhook()
	payload := "{}"

	timestamp := time.Now().Unix() - (DefaultTolerance*60 + 1)
	signature := generateTestSignature(payload, timestamp, testSecretKey)

	_, err := webhook.ConstructEvent(payload, signature, DefaultTolerance)

	if err == nil {
		t.Error("Expected an error for timestamp outside tolerance")
		return
	}

	invalidSigErr, ok := err.(*InvalidSignatureError)
	if !ok {
		t.Errorf("Expected InvalidSignatureError, got %T", err)
		return
	}

	if invalidSigErr.Message != "Timestamp is outside the tolerance zone." {
		t.Errorf("Expected 'Timestamp is outside the tolerance zone.', got '%s'", invalidSigErr.Message)
	}
}

func TestZeroTolerance(t *testing.T) {
	webhook := getTestWebhook()
	payload := `{"version":1,"type":"test","id":"1","source":"test","time":"2025-01-01T00:00:00Z","tenantId":"t","data":{}}`

	timestamp := int64(1630000000)
	signature := generateTestSignature(payload, timestamp, testSecretKey)

	event, err := webhook.ConstructEvent(payload, signature, 0)

	if err != nil {
		t.Errorf("Expected no error with zero tolerance, got %v", err)
		return
	}

	if event == nil {
		t.Error("Expected event to be returned")
	}
}

func TestSignatureWithEqualsInValue(t *testing.T) {
	webhook := getTestWebhook()
	payload := `{"version":1,"type":"test","id":"1","source":"test","time":"2025-01-01T00:00:00Z","tenantId":"t","data":{}}`

	timestamp := time.Now().Unix()
	signature := generateTestSignature(payload, timestamp, testSecretKey)

	event, err := webhook.ConstructEvent(payload, signature, -1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if event == nil {
		t.Error("Expected event to be returned")
	}
}
