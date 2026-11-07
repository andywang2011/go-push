package push

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewAPNsInvalidConfig(t *testing.T) {
	_, err := NewAPNs(APNsConfig{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAPNsSendSuccess(t *testing.T) {
	var topic, auth string
	var payload map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth = r.Header.Get("Authorization")
		topic = r.Header.Get("apns-topic")
		if !strings.HasPrefix(r.URL.Path, "/3/device/") {
			t.Errorf("path = %q", r.URL.Path)
		}
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &payload)
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	sender, err := NewAPNs(APNsConfig{
		KeyID:         "KEY",
		TeamID:        "TEAM",
		BundleID:      "com.example.app",
		TokenProvider: staticToken("signed-jwt"),
		BaseURL:       srv.URL,
		HTTPClient:    srv.Client(),
	})
	if err != nil {
		t.Fatal(err)
	}

	err = sender.Send(context.Background(), Message{
		Tokens: []string{"abc123"},
		Title:  "Alert",
		Body:   "Message",
		Data:   map[string]string{"deeplink": "/home"},
	})
	if err != nil {
		t.Fatalf("Send() = %v", err)
	}
	if auth != "bearer signed-jwt" {
		t.Errorf("Authorization = %q", auth)
	}
	if topic != "com.example.app" {
		t.Errorf("apns-topic = %q", topic)
	}
	aps, _ := payload["aps"].(map[string]interface{})
	if aps == nil {
		t.Fatalf("payload = %v", payload)
	}
}

func TestAPNsSendHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusGone)
		_, _ = io.WriteString(w, `{"reason":"Unregistered"}`)
	}))
	t.Cleanup(srv.Close)

	sender, err := NewAPNs(APNsConfig{
		KeyID: "K", TeamID: "T", BundleID: "com.app",
		TokenProvider: staticToken("jwt"),
		BaseURL:       srv.URL,
		HTTPClient:    srv.Client(),
	})
	if err != nil {
		t.Fatal(err)
	}

	err = sender.Send(context.Background(), Message{Tokens: []string{"dead"}, Body: "x"})
	if err == nil || !strings.Contains(err.Error(), "410") {
		t.Fatalf("Send() = %v", err)
	}
}
