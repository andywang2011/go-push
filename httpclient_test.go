package push

import (
	"net/http"
	"testing"
	"time"
)

func TestDefaultHTTPClient_HasTimeout(t *testing.T) {
	c := defaultHTTPClient()
	if c == http.DefaultClient {
		t.Fatal("must not return http.DefaultClient")
	}
	if c.Timeout != 10*time.Second {
		t.Fatalf("Timeout = %v, want 10s", c.Timeout)
	}
	if defaultHTTPClient() == c {
		t.Fatal("each call must return a fresh client")
	}
}

func TestNew_DefaultHTTPClientHasTimeout(t *testing.T) {
	fcm, err := NewFCM(FCMConfig{ProjectID: "p", TokenProvider: staticToken("t")})
	if err != nil {
		t.Fatal(err)
	}
	assertTimeout(t, "fcm", fcm.(*fcmSender).client)

	apns, err := NewAPNs(APNsConfig{KeyID: "K", TeamID: "T", BundleID: "b", TokenProvider: staticToken("t")})
	if err != nil {
		t.Fatal(err)
	}
	assertTimeout(t, "apns", apns.(*apnsSender).client)
}

func TestNew_CustomHTTPClientPreserved(t *testing.T) {
	custom := &http.Client{Timeout: time.Minute}
	fcm, err := NewFCM(FCMConfig{ProjectID: "p", TokenProvider: staticToken("t"), HTTPClient: custom})
	if err != nil {
		t.Fatal(err)
	}
	if fcm.(*fcmSender).client != HTTPDoer(custom) {
		t.Fatal("custom HTTPClient must be used as is")
	}
}

func assertTimeout(t *testing.T, name string, d HTTPDoer) {
	t.Helper()
	c, ok := d.(*http.Client)
	if !ok || c == http.DefaultClient || c.Timeout != 10*time.Second {
		t.Fatalf("%s default client = %#v, want *http.Client with 10s timeout", name, d)
	}
}
