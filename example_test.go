package push_test

import (
	"bytes"
	"context"
	"fmt"

	"github.com/andywang2011/go-push"
)

func ExampleNewLogSender() {
	buf := &bytes.Buffer{}
	dev := push.NewLogSender(buf)
	err := dev.Send(context.Background(), push.Message{
		Tokens: []string{"device-token"},
		Title:  "Hello",
		Body:   "World",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.Len() > 0)
	// Output: true
}
