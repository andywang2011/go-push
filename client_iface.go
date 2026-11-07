package push

import (
	"context"
	"net/http"
)

// Sender delivers push notifications (mock-friendly client surface).
type Sender interface {
	// Send delivers msg using the provider API or sink.
	Send(ctx context.Context, msg Message) error
}

// HTTPDoer performs HTTP requests; *http.Client implements this interface.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// AccessTokenProvider supplies OAuth2 bearer tokens (FCM HTTP v1).
type AccessTokenProvider interface {
	// AccessToken returns a short-lived access token for the FCM API.
	AccessToken(ctx context.Context) (string, error)
}
