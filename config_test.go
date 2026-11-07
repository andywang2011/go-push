package push

import "testing"

func TestFCMConfigValidate(t *testing.T) {
	if err := (FCMConfig{}).Validate(); err == nil {
		t.Fatal("expected error")
	}
	cfg := FCMConfig{ProjectID: "p", TokenProvider: staticToken("t")}
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestAPNsConfigValidate(t *testing.T) {
	if err := (APNsConfig{}).Validate(); err == nil {
		t.Fatal("expected error")
	}
	cfg := APNsConfig{
		KeyID: "K", TeamID: "T", BundleID: "com.app",
		TokenProvider: staticToken("jwt"),
	}
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestAPNsResolvedBaseURL(t *testing.T) {
	sandbox := APNsConfig{Sandbox: true}
	if sandbox.resolvedBaseURL() != defaultAPNsSandboxURL {
		t.Fatalf("sandbox url = %q", sandbox.resolvedBaseURL())
	}
	prod := APNsConfig{Sandbox: false}
	if prod.resolvedBaseURL() != defaultAPNsProductionURL {
		t.Fatalf("prod url = %q", prod.resolvedBaseURL())
	}
}
