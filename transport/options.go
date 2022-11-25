package transport

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

// HookRequest hook request before http request function
type HookRequest func(req *http.Request) error

type Options struct {
	Headers       map[string]string // http request header
	Client        *http.Client      // http client
	MaxRetries    int               // max retry if error
	Timeout       time.Duration     // http connect timeout
	RetryInterval time.Duration     // retry interval if error
	Trans         *http.Transport   // http transport
	Hook          HookRequest       // hook before http request
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

// Headers http header, if there are more than two keys, use the newest one
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

// Client custom http client
func Client(c *http.Client) Option {
	return func(o *Options) { o.Client = c }
}

// MaxRetries max retry times
func MaxRetries(r int) Option {
	return func(o *Options) { o.MaxRetries = r }
}

// Timeout connect timeout
func Timeout(t time.Duration) Option {
	return func(o *Options) { o.Timeout = t }
}

// RetryInterval retry interval if error
func RetryInterval(t time.Duration) Option {
	return func(o *Options) { o.RetryInterval = t }
}

// Trans http transport
func Trans(t *http.Transport) Option {
	return func(o *Options) { o.Trans = t }
}

func Hook(h HookRequest) Option {
	return func(o *Options) { o.Hook = h }
}
