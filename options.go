package apollo

import (
	"github.com/xnzone/apollo-go/auth"
	"github.com/xnzone/apollo-go/transport"
)

type Option func(o *Options)

type Options struct {
	Auth      auth.Auth
	Transport transport.Transport
}

func newOptions(opts ...Option) *Options {
	opt := Options{
		Auth:      auth.DefaultAuth,
		Transport: transport.DefaultTransport,
	}
	for _, o := range opts {
		o(&opt)
	}
	return &opt
}
