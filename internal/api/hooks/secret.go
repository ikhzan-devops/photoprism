package hooks

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/photoprism/photoprism/pkg/media"
	"github.com/photoprism/photoprism/pkg/media/http/header"
)

var timeTolerance = 5 * time.Minute

var (
	ErrRequiredHeaders     = errors.New("missing required headers")
	ErrInvalidHeaders      = errors.New("invalid signature headers")
	ErrNoMatchingSignature = errors.New("no matching signature found")
	ErrMessageTooOld       = errors.New("message timestamp too old")
	ErrMessageTooNew       = errors.New("message timestamp too new")
)

type Secret struct {
	key []byte
}

func NewSecret(secret string) (*Secret, error) {
	key, err := media.DecodeBase64String(strings.TrimPrefix(secret, header.WebhookSecretPrefix))
	if err != nil {
		return nil, fmt.Errorf("unable to create webhook, err: %w", err)
	}
	return &Secret{
		key: key,
	}, nil
}

func NewWebhookRaw(secret []byte) (*Secret, error) {
	return &Secret{
		key: secret,
	}, nil
}

// Verify validates the payload against the webhook signature headers
// using the webhooks signing secret.
//
// Returns an error if the body or headers are missing/unreadable
// or if the signature doesn't match.
func (wh *Secret) Verify(payload []byte, headers http.Header) error {
	return wh.verify(payload, headers, true)
}

// VerifyIgnoringTimestamp validates the payload against the webhook signature headers
// using the webhooks signing secret.
//
// Returns an error if the body or headers are missing/unreadable
// or if the signature doesn't match.
//
// WARNING: This function does not check the signature's timestamp.
// We recommend using the `Verify` function instead.
func (wh *Secret) VerifyIgnoringTimestamp(payload []byte, headers http.Header) error {
	return wh.verify(payload, headers, false)
}

func (wh *Secret) verify(payload []byte, headers http.Header, enforceTolerance bool) error {
	msgId := headers.Get(header.WebhookID)
	msgSignature := headers.Get(header.WebhookSignature)
	msgTimestamp := headers.Get(header.WebhookTimestamp)
	if msgId == "" || msgSignature == "" || msgTimestamp == "" {
		return fmt.Errorf("unable to verify payload, err: %w", ErrRequiredHeaders)
	}

	timestamp, err := parseTimestampHeader(msgTimestamp)
	if err != nil {
		return fmt.Errorf("unable to verify payload, err: %w", err)
	}

	if enforceTolerance {
		if err := verifyTimestamp(timestamp); err != nil {
			return fmt.Errorf("unable to verify payload, err: %w", err)
		}
	}

	_, expectedSignature, err := wh.sign(msgId, timestamp, payload)
	if err != nil {
		return fmt.Errorf("unable to verify payload, err: %w", err)
	}

	passedSignatures := strings.Split(msgSignature, " ")
	for _, versionedSignature := range passedSignatures {
		sigParts := strings.Split(versionedSignature, ",")
		if len(sigParts) < 2 {
			continue
		}

		version := sigParts[0]

		if version != "v1" {
			continue
		}

		signature := []byte(sigParts[1])

		if hmac.Equal(signature, expectedSignature) {
			return nil
		}
	}

	return fmt.Errorf("unable to verify payload, err: %w", ErrNoMatchingSignature)
}

func (wh *Secret) Sign(msgId string, timestamp time.Time, payload []byte) (string, error) {
	version, signature, err := wh.sign(msgId, timestamp, payload)
	return fmt.Sprintf("%s,%s", version, signature), err
}

func (wh *Secret) sign(msgId string, timestamp time.Time, payload []byte) (version string, signature []byte, err error) {
	toSign := fmt.Sprintf("%s.%d.%s", msgId, timestamp.Unix(), payload)

	h := hmac.New(sha256.New, wh.key)
	h.Write([]byte(toSign))
	sig := make([]byte, media.EncodedLenBase64(h.Size()))
	media.EncodeBase64Bytes(sig, h.Sum(nil))

	return "v1", sig, nil
}

func parseTimestampHeader(timestampHeader string) (time.Time, error) {
	timeInt, err := strconv.ParseInt(timestampHeader, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("unable to parse timestamp header, err: %w", errors.Join(err, ErrInvalidHeaders))
	}
	timestamp := time.Unix(timeInt, 0)
	return timestamp, nil
}

func verifyTimestamp(timestamp time.Time) error {
	now := time.Now()

	if now.Sub(timestamp) > timeTolerance {
		return ErrMessageTooOld
	}

	if timestamp.After(now.Add(timeTolerance)) {
		return ErrMessageTooNew
	}

	return nil
}
