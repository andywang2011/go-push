// Package push provides multi-provider mobile push notification sending with config injection.
//
// Design constraints:
//   - No viper, environment variables, or config file paths inside the library.
//   - No global singletons or log.Fatal; callers construct senders via New*.
//   - FCM and APNs accept an optional HTTP client and token providers in Config for httptest.
//   - New* returns Sender (client_iface.go) for mocks and dependency injection.
//   - ConfigDefault: N/A — use FCMConfig / APNsConfig per provider.
package push
