package auth

type Auth interface {
	Headers(uri string, appID string, secret string) map[string][]string
}
