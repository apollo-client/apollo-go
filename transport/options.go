package transport

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

type HookRequest func(req *http.Request) error

type Options struct {
	Headers       map[string]string
	Client        *http.Client
	MaxRetries    int
	Timeout       time.Duration
	RetryInterval time.Duration
	Trans         *http.Transport
	Hook          HookRequest
}

var (
	trans = &http.Transport{
		MaxIdleConns:        512,
		MaxIdleConnsPerHost: 512,
		DialContext: (&net.Dialer{
			KeepAlive: 60 * time.Second,
			Timeout:   1 * time.Second,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
)

func newOptions(opts ...Option) *Options {
	opt := Options{
		Client:        &http.Client{},
		Timeout:       1 * time.Second,
		RetryInterval: 2 * time.Second,
		MaxRetries:    5,
		Trans:         trans,
	}
	initOptions(&opt, opts...)
	return &opt
}

func initOptions(opt *Options, opts ...Option) {
	for _, o := range opts {
		o(opt)
	}
}

func Headers(m map[string]string) Option {
	return func(o *Options) {
		if o.Headers == nil {
			o.Headers = make(map[string]string)
		}
		for k, v := range m {
			o.Headers[k] = v
		}
	}
}

func Client(c *http.Client) Option {
	return func(o *Options) { o.Client = c }
}

func MaxRetries(r int) Option {
	return func(o *Options) { o.MaxRetries = r }
}

func Timeout(t time.Duration) Option {
	return func(o *Options) { o.Timeout = t }
}

func RetryInterval(t time.Duration) Option {
	return func(o *Options) { o.RetryInterval = t }
}

func Trans(t *http.Transport) Option {
	return func(o *Options) { o.Trans = t }
}

func Hook(h HookRequest) Option {
	return func(o *Options) { o.Hook = h }
}