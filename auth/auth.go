package auth

type Auth interface {
	Header(rawURL string, appId string, secret string) map[string]string
}

var (
	DefaultAuth = newAuth()
)

func newAuth() Auth {
	return NewApolloAuth()
}
