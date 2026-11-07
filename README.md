# go-push

Multi-provider mobile push notifications for Go (FCM, APNs, log sender).

## Install

```bash
go get github.com/andywang2011/go-push@v0.1.0
```

## Quick start

```go
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
```

## Examples

```bash
cd examples
go run .
```

## License

MIT — see [LICENSE](LICENSE).
