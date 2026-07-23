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

const DefaultTolerance = 5

const version = "v2"

type Webhook struct {
	ApiSecretKey string
}

func NewWebhook(apiSecretKey string) *Webhook {
	return &Webhook{ApiSecretKey: apiSecretKey}
}

func (w *Webhook) ConstructEvent(payload string, signature string, tolerance int) (*WebhookEvent, error) {
	if err := w.verifySignature(payload, signature, tolerance); err != nil {
		return nil, err
	}

	var envelope map[string]json.RawMessage
	if err := json.Unmarshal([]byte(payload), &envelope); err != nil || envelope == nil {
		return nil, NewInvalidPayloadError("Payload format is invalid.")
	}
	if _, isBatch := envelope["records"]; isBatch {
		return nil, NewInvalidPayloadError("Payload is a batch of log events. Use ConstructLogEventBatch instead.")
	}

	var event WebhookEvent
	if err := json.Unmarshal([]byte(payload), &event); err != nil {
		return nil, NewInvalidPayloadError("Payload format is invalid.")
	}
	if err := validateWebhookEvent(&event); err != nil {
		return nil, err
	}

	return &event, nil
}

func (w *Webhook) ConstructEventWithDefaultTolerance(payload string, signature string) (*WebhookEvent, error) {
	return w.ConstructEvent(payload, signature, DefaultTolerance)
}

func (w *Webhook) ConstructLogEventBatch(payload string, signature string, tolerance int) (*WebhookEventBatch, error) {
	if err := w.verifySignature(payload, signature, tolerance); err != nil {
		return nil, err
	}

	var batch WebhookEventBatch
	if err := json.Unmarshal([]byte(payload), &batch); err != nil || batch.Records == nil {
		return nil, NewInvalidPayloadError("Payload format is invalid. Expected a 'records' array.")
	}
	for i := range batch.Records {
		if err := validateWebhookLogEvent(&batch.Records[i]); err != nil {
			return nil, err
		}
	}

	return &batch, nil
}

func (w *Webhook) ConstructLogEventBatchWithDefaultTolerance(payload string, signature string) (*WebhookEventBatch, error) {
	return w.ConstructLogEventBatch(payload, signature, DefaultTolerance)
}

func (w *Webhook) verifySignature(payload string, signature string, tolerance int) error {
	parsedSignature, err := w.parseSignature(signature)
	if err != nil {
		return err
	}

	secondsSinceEpoch := time.Now().Unix()

	if tolerance > 0 && parsedSignature.Timestamp < secondsSinceEpoch-int64(tolerance*60) {
		return NewInvalidSignatureError("Timestamp is outside the tolerance zone.")
	}

	hmacContent := strconv.FormatInt(parsedSignature.Timestamp, 10) + "." + payload

	computedSignature := w.computeHmac(hmacContent)

	match := false
	for _, sig := range parsedSignature.Signatures {
		if hmac.Equal([]byte(sig), []byte(computedSignature)) {
			match = true
			break
		}
	}

	if !match {
		return NewInvalidSignatureError("Signature mismatch.")
	}

	return nil
}

func validateWebhookEvent(event *WebhookEvent) error {
	if err := validateWebhookEnvelope(
		event.Version,
		event.Type,
		event.Id,
		event.Source,
		event.Time,
		event.TenantId,
	); err != nil {
		return err
	}
	if event.Data == nil {
		return NewInvalidPayloadError("Payload is missing required field 'data'.")
	}
	return nil
}

func validateWebhookLogEvent(event *WebhookLogEvent) error {
	if err := validateWebhookEnvelope(
		event.Version,
		event.Type,
		event.Id,
		event.Source,
		event.Time,
		event.TenantId,
	); err != nil {
		return err
	}
	if event.Record == nil {
		return NewInvalidPayloadError("Payload is missing required field 'record'.")
	}
	return nil
}

func validateWebhookEnvelope(version int, eventType string, id string, source string, eventTime string, tenantId string) error {
	if version <= 0 {
		return NewInvalidPayloadError("Payload is missing required field 'version'.")
	}

	requiredFields := []struct {
		name  string
		value string
	}{
		{name: "type", value: eventType},
		{name: "id", value: id},
		{name: "source", value: source},
		{name: "time", value: eventTime},
		{name: "tenantId", value: tenantId},
	}
	for _, field := range requiredFields {
		if strings.TrimSpace(field.value) == "" {
			return NewInvalidPayloadError("Payload is missing required field '" + field.name + "'.")
		}
	}

	return nil
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
	return strings.Replace(signature, "=", "", -1)
}
