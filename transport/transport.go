package transport

type Option func(o *Options)
type CallOption func(o *CallOptions)

// Transport request transport
type Transport interface {
	Init(opt ...Option) error
	Options() *Options
	Do(reqURL string, opt ...CallOption) (status int, body []byte, err error)
}

var (
	DefaultTransport = newTransport()
)

func newTransport() Transport {
	return NewHTTPTransport()
}
