package push

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestLogSenderSuccess(t *testing.T) {
	var buf bytes.Buffer
	s := NewLogSender(&buf)
	err := s.Send(context.Background(), Message{
		Tokens: []string{"tok1"},
		Title:  "Hi",
		Body:   "There",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "tok1") {
		t.Fatalf("output = %q", buf.String())
	}
}

func TestLogSenderValidation(t *testing.T) {
	s := NewLogSender(nil)
	if err := s.Send(context.Background(), Message{}); err == nil {
		t.Fatal("expected validation error")
	}
}
