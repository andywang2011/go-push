package push

import (
	"errors"
	"strings"
)

const (
	defaultFCMBaseURL = "https://fcm.googleapis.com"
	// Doc: https://firebase.google.com/docs/reference/fcm/rest/v1/projects.messages/send
	fcmSendPathPattern = "/v1/projects/%s/messages:send"

	defaultAPNsProductionURL = "https://api.push.apple.com"
	defaultAPNsSandboxURL    = "https://api.sandbox.push.apple.com"
	// Doc: https://developer.apple.com/documentation/usernotifications/sending-notification-requests-to-apns
	apnsDevicePathPattern = "/3/device/%s"
)

// FCMConfig holds Firebase Cloud Messaging HTTP v1 credentials and transport options.
type FCMConfig struct {
	// ProjectID is the Firebase/Google Cloud project ID (required).
	ProjectID string
	// TokenProvider supplies OAuth2 access tokens for the FCM API (required).
	TokenProvider AccessTokenProvider
	// BaseURL overrides the FCM API base URL (for httptest).
	BaseURL string
	// HTTPClient is the HTTP client used for requests; nil uses a client with a 10s timeout.
	HTTPClient HTTPDoer
}

// Validate checks required FCM fields.
func (c FCMConfig) Validate() error {
	if strings.TrimSpace(c.ProjectID) == "" {
		return errors.New("fcm project id is not set")
	}
	if c.TokenProvider == nil {
		return errors.New("fcm token provider is not set")
	}
	return nil
}

func (c FCMConfig) resolvedBaseURL() string {
	if strings.TrimSpace(c.BaseURL) != "" {
		return strings.TrimRight(c.BaseURL, "/")
	}
	return defaultFCMBaseURL
}

// APNsConfig holds Apple Push Notification service token-based (.p8) settings.
type APNsConfig struct {
	// KeyID is the APNs auth key identifier (required).
	KeyID string
	// TeamID is the Apple Developer Team ID (required).
	TeamID string
	// BundleID is the app's bundle identifier used as apns-topic (required).
	BundleID string
	// TokenProvider supplies JWT bearer tokens for APNs (required).
	TokenProvider AccessTokenProvider
	// Sandbox uses the APNs sandbox endpoint when true.
	Sandbox bool
	// BaseURL overrides the APNs API base URL (for httptest).
	BaseURL string
	// HTTPClient is the HTTP client used for requests; nil uses a client with a 10s timeout.
	HTTPClient HTTPDoer
}

// Validate checks required APNs fields.
func (c APNsConfig) Validate() error {
	if strings.TrimSpace(c.KeyID) == "" {
		return errors.New("apns key id is not set")
	}
	if strings.TrimSpace(c.TeamID) == "" {
		return errors.New("apns team id is not set")
	}
	if strings.TrimSpace(c.BundleID) == "" {
		return errors.New("apns bundle id is not set")
	}
	if c.TokenProvider == nil {
		return errors.New("apns token provider is not set")
	}
	return nil
}

func (c APNsConfig) resolvedBaseURL() string {
	if strings.TrimSpace(c.BaseURL) != "" {
		return strings.TrimRight(c.BaseURL, "/")
	}
	if c.Sandbox {
		return defaultAPNsSandboxURL
	}
	return defaultAPNsProductionURL
}
