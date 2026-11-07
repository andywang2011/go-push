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

type fcmSender struct {
	projectID     string
	baseURL       string
	tokenProvider AccessTokenProvider
	client        HTTPDoer
}

// NewFCM creates a Firebase Cloud Messaging HTTP v1 sender from cfg.
func NewFCM(cfg FCMConfig) (Sender, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	client := cfg.HTTPClient
	if client == nil {
		client = defaultHTTPClient()
	}
	return &fcmSender{
		projectID:     cfg.ProjectID,
		baseURL:       cfg.resolvedBaseURL(),
		tokenProvider: cfg.TokenProvider,
		client:        client,
	}, nil
}

// Send posts one FCM v1 message per device token.
func (s *fcmSender) Send(ctx context.Context, msg Message) error {
	if err := validateMessage(msg); err != nil {
		return err
	}

	token, err := s.tokenProvider.AccessToken(ctx)
	if err != nil {
		return fmt.Errorf("fcm access token: %w", err)
	}
	if strings.TrimSpace(token) == "" {
		return fmt.Errorf("fcm access token is empty")
	}

	for _, deviceToken := range msg.Tokens {
		if err := s.sendOne(ctx, token, deviceToken, msg); err != nil {
			return err
		}
	}
	return nil
}

func (s *fcmSender) sendOne(ctx context.Context, accessToken, deviceToken string, msg Message) error {
	payload := map[string]interface{}{
		"message": buildFCMMessage(deviceToken, msg),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("fcm encode request: %w", err)
	}

	path := fmt.Sprintf(fcmSendPathPattern, s.projectID)
	reqURL := s.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("fcm create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("fcm request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("fcm api error: status %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}
	return nil
}

func buildFCMMessage(token string, msg Message) map[string]interface{} {
	out := map[string]interface{}{
		"token": token,
	}
	if msg.Title != "" || msg.Body != "" {
		out["notification"] = map[string]string{
			"title": msg.Title,
			"body":  msg.Body,
		}
	}
	if len(msg.Data) > 0 {
		out["data"] = msg.Data
	}
	if msg.Platform == PlatformAndroid {
		out["android"] = map[string]interface{}{"priority": "HIGH"}
	}
	return out
}
