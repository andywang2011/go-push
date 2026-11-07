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

type staticToken string

func (s staticToken) AccessToken(context.Context) (string, error) {
	return string(s), nil
}

func TestNewFCMInvalidConfig(t *testing.T) {
	_, err := NewFCM(FCMConfig{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFCMSendSuccess(t *testing.T) {
	var gotAuth string
	var gotBody map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		if !strings.HasSuffix(r.URL.Path, "/messages:send") {
			t.Errorf("path = %q", r.URL.Path)
		}
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &gotBody)
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"name":"projects/p/messages/0"}`)
	}))
	t.Cleanup(srv.Close)

	sender, err := NewFCM(FCMConfig{
		ProjectID:     "my-project",
		TokenProvider: staticToken("oauth-token"),
		BaseURL:       srv.URL,
		HTTPClient:    srv.Client(),
	})
	if err != nil {
		t.Fatal(err)
	}

	err = sender.Send(context.Background(), Message{
		Tokens:   []string{"device-a", "device-b"},
		Title:    "Hello",
		Body:     "World",
		Data:     map[string]string{"k": "v"},
		Platform: PlatformAndroid,
	})
	if err != nil {
		t.Fatalf("Send() = %v", err)
	}
	if gotAuth != "Bearer oauth-token" {
		t.Errorf("Authorization = %q", gotAuth)
	}
	msg, _ := gotBody["message"].(map[string]interface{})
	if msg == nil || msg["token"] != "device-b" {
		t.Errorf("last message token = %v", msg)
	}
}

func TestFCMSendHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, `{"error":"bad token"}`)
	}))
	t.Cleanup(srv.Close)

	sender, err := NewFCM(FCMConfig{
		ProjectID:     "p",
		TokenProvider: staticToken("t"),
		BaseURL:       srv.URL,
		HTTPClient:    srv.Client(),
	})
	if err != nil {
		t.Fatal(err)
	}

	err = sender.Send(context.Background(), Message{Tokens: []string{"x"}, Title: "t"})
	if err == nil || !strings.Contains(err.Error(), "400") {
		t.Fatalf("Send() = %v", err)
	}
}

func TestFCMEmptyAccessToken(t *testing.T) {
	sender, err := NewFCM(FCMConfig{
		ProjectID:     "p",
		TokenProvider: staticToken("  "),
	})
	if err != nil {
		t.Fatal(err)
	}
	err = sender.Send(context.Background(), Message{Tokens: []string{"x"}, Body: "b"})
	if err == nil || !strings.Contains(err.Error(), "empty") {
		t.Fatalf("Send() = %v", err)
	}
}
