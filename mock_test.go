package push

import (
	"context"
	"errors"
	"testing"
)

type fakeSender struct {
	sendFn func(context.Context, Message) error
}

func (f fakeSender) Send(ctx context.Context, msg Message) error {
	if f.sendFn == nil {
		return nil
	}
	return f.sendFn(ctx, msg)
}

func TestFakeSenderPattern(t *testing.T) {
	var called bool
	s := fakeSender{sendFn: func(_ context.Context, msg Message) error {
		called = true
		if len(msg.Tokens) == 0 {
			return errors.New("empty")
		}
		return nil
	}}
	if err := s.Send(context.Background(), Message{Tokens: []string{"t"}}); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("send not called")
	}
}
