package transport

type Option func(o *Options)

// Transport request transport
type Transport interface {
	Init(opt ...Option) error
	Options() *Options
	Do(reqURL string, opt ...Option) (status int, body []byte, err error)
}

var (
	DefaultTransport = newTransport()
)

func newTransport() Transport {
	return NewHTTPTransport()
}
