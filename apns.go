package push

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type apnsSender struct {
	bundleID      string
	baseURL       string
	tokenProvider AccessTokenProvider
	client        HTTPDoer
}

// NewAPNs creates an Apple Push Notification service sender from cfg.
func NewAPNs(cfg APNsConfig) (Sender, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	client := cfg.HTTPClient
	if client == nil {
		client = defaultHTTPClient()
	}
	return &apnsSender{
		bundleID:      cfg.BundleID,
		baseURL:       cfg.resolvedBaseURL(),
		tokenProvider: cfg.TokenProvider,
		client:        client,
	}, nil
}

// Send posts one notification request per device token.
func (s *apnsSender) Send(ctx context.Context, msg Message) error {
	if err := validateMessage(msg); err != nil {
		return err
	}

	jwt, err := s.tokenProvider.AccessToken(ctx)
	if err != nil {
		return fmt.Errorf("apns access token: %w", err)
	}
	if strings.TrimSpace(jwt) == "" {
		return fmt.Errorf("apns jwt is empty")
	}

	for _, deviceToken := range msg.Tokens {
		if err := s.sendOne(ctx, jwt, deviceToken, msg); err != nil {
			return err
		}
	}
	return nil
}

func (s *apnsSender) sendOne(ctx context.Context, jwt, deviceToken string, msg Message) error {
	body, err := json.Marshal(buildAPNsPayload(msg))
	if err != nil {
		return fmt.Errorf("apns encode request: %w", err)
	}

	path := fmt.Sprintf(apnsDevicePathPattern, deviceToken)
	reqURL := s.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("apns create request: %w", err)
	}
	req.Header.Set("Authorization", "bearer "+jwt)
	req.Header.Set("apns-topic", s.bundleID)
	req.Header.Set("apns-push-type", "alert")
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("apns request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("apns api error: status %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}
	return nil
}

func buildAPNsPayload(msg Message) map[string]interface{} {
	aps := map[string]interface{}{}
	if msg.Title != "" || msg.Body != "" {
		aps["alert"] = map[string]string{
			"title": msg.Title,
			"body":  msg.Body,
		}
	}
	aps["sound"] = "default"

	out := map[string]interface{}{
		"aps": aps,
	}
	for k, v := range msg.Data {
		out[k] = v
	}
	return out
}
