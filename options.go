package apollo

import (
	"time"

	"github.com/xnzone/apollo-go/auth"
	"github.com/xnzone/apollo-go/transport"
)

type Option func(o *Options)

type Options struct {
	Auth          auth.Auth           // auth interface
	Transport     transport.Transport // transport interface
	WatchInterval time.Duration       // watch interval
}

func newOptions(opts ...Option) *Options {
	opt := Options{
		Auth:          auth.DefaultAuth,
		Transport:     transport.DefaultTransport,
		WatchInterval: 5 * time.Second,
	}
	for _, o := range opts {
		o(&opt)
	}
	return &opt
}

func Auth(a auth.Auth) Option {
	return func(o *Options) { o.Auth = a }
}

func Transport(t transport.Transport) Option {
	return func(o *Options) { o.Transport = t }
}

func WatchInterval(t time.Duration) Option {
	return func(o *Options) { o.WatchInterval = t }
}
