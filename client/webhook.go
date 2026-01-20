package client

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// Default tolerance (in minutes) for difference between timestamp in signature and current time
// This is used to prevent replay attacks
const DefaultTolerance = 5

const version = "v2"

type Webhook struct {
	ApiSecretKey string
}

func NewWebhook(apiSecretKey string) *Webhook {
	return &Webhook{ApiSecretKey: apiSecretKey}
}

// ConstructEvent verifies the webhook signature and returns the parsed event.
// payload is the raw request body as a string.
// signature is the value of the webhook signature header.
// tolerance is the maximum age of the webhook in minutes (use DefaultTolerance or -1 to disable).
func (w *Webhook) ConstructEvent(payload string, signature string, tolerance int) (*WebhookEvent, error) {
	parsedSignature, err := w.parseSignature(signature)
	if err != nil {
		return nil, err
	}

	secondsSinceEpoch := time.Now().Unix()

	if tolerance > 0 && parsedSignature.Timestamp < secondsSinceEpoch-int64(tolerance*60) {
		return nil, NewInvalidSignatureError("Timestamp is outside the tolerance zone.")
	}

	hmacContent := strconv.FormatInt(parsedSignature.Timestamp, 10) + "." + payload

	computedSignature := w.computeHmac(hmacContent)

	match := false
	for _, sig := range parsedSignature.Signatures {
		// Use constant-time comparison to prevent timing attacks
		if hmac.Equal([]byte(sig), []byte(computedSignature)) {
			match = true
			break
		}
	}

	if !match {
		return nil, NewInvalidSignatureError("Signature mismatch.")
	}

	var event WebhookEvent
	if err := json.Unmarshal([]byte(payload), &event); err != nil {
		return nil, err
	}

	return &event, nil
}

// ConstructEventWithDefaultTolerance verifies the webhook signature using the default tolerance.
func (w *Webhook) ConstructEventWithDefaultTolerance(payload string, signature string) (*WebhookEvent, error) {
	return w.ConstructEvent(payload, signature, DefaultTolerance)
}

type signatureHeaderData struct {
	Signatures []string
	Timestamp  int64
}

func (w *Webhook) parseSignature(value string) (*signatureHeaderData, error) {
	if value == "" {
		return nil, NewInvalidSignatureError("Signature format is invalid.")
	}

	result := &signatureHeaderData{
		Timestamp:  -1,
		Signatures: []string{},
	}

	items := strings.Split(value, ",")
	for _, item := range items {
		kv := strings.SplitN(item, "=", 2)
		if len(kv) != 2 {
			continue
		}

		if kv[0] == "t" {
			timestamp, err := strconv.ParseInt(kv[1], 10, 64)
			if err != nil {
				continue
			}
			result.Timestamp = timestamp
		}

		if kv[0] == version {
			result.Signatures = append(result.Signatures, kv[1])
		}
	}

	if result.Timestamp == -1 || len(result.Signatures) == 0 {
		return nil, NewInvalidSignatureError("Signature format is invalid.")
	}

	return result, nil
}

func (w *Webhook) computeHmac(data string) string {
	mac := hmac.New(sha256.New, []byte(w.ApiSecretKey))
	mac.Write([]byte(data))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	// Remove trailing '=' characters to match the expected format
	return strings.Replace(signature, "=", "", -1)
}
