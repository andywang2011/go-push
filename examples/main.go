package main

import (
	"context"
	"os"

	"github.com/andywang2011/go-push"
)

func main() {
	dev := push.NewLogSender(os.Stdout)
	_ = dev.Send(context.Background(), push.Message{
		Tokens: []string{"device-token"},
		Title:  "Hello",
		Body:   "From go-push examples",
	})
}
