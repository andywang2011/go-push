package push

import (
	"errors"
	"strings"
)

const (
	// PlatformAndroid is the Android / FCM platform identifier.
	PlatformAndroid = "android"
	// PlatformIOS is the Apple / APNs platform identifier.
	PlatformIOS = "ios"
)

// Message is the outbound push payload shared by all providers.
type Message struct {
	// Tokens are device registration tokens (FCM) or APNs device tokens.
	Tokens []string
	// Title is the notification title shown to the user.
	Title string
	// Body is the notification body text.
	Body string
	// Data is custom key-value payload delivered to the app.
	Data map[string]string
	// Platform hints the target stack: "android", "ios", or empty for provider default.
	Platform string
}

func validateMessage(msg Message) error {
	if len(msg.Tokens) == 0 {
		return errors.New("push tokens are not set")
	}
	for i, tok := range msg.Tokens {
		if strings.TrimSpace(tok) == "" {
			return errors.New("push token is empty")
		}
		msg.Tokens[i] = strings.TrimSpace(tok)
	}
	if strings.TrimSpace(msg.Title) == "" && strings.TrimSpace(msg.Body) == "" && len(msg.Data) == 0 {
		return errors.New("push message must have title, body, or data")
	}
	return nil
}
