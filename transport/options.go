package transport

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

type Options struct {
	Client    *http.Client
	Retries   int32
	Trans     *http.Transport
	Insecure  bool
	TLSConfig *tls.Config
	Timeout   time.Duration
}

var (
	trans = &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		MaxIdleConns:        512,
		MaxIdleConnsPerHost: 512,
		DialContext: (&net.Dialer{
			KeepAlive: 60 * time.Second,
			Timeout:   1 * time.Second,
		}).DialContext,
	}
)

func newOptions(opt ...Option) Options {
	opts := Options{
		Client:  &http.Client{},
		Retries: 5,
		Trans:   trans,
		Timeout: 2 * time.Second,
	}
	initOptions(&opts, opt...)

	return opts
}

func initOptions(opts *Options, opt ...Option) {
	for _, o := range opt {
		o(opts)
	}
	if opts.Insecure {
		opts.Trans.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	if opts.TLSConfig != nil {
		opts.Trans.TLSClientConfig = opts.TLSConfig
	}
}

func Timeout(t time.Duration) Option {
	return func(o *Options) { o.Timeout = t }
}
func Client(c *http.Client) Option {
	return func(o *Options) {
		o.Client = c
	}
}

func Retries(r int32) Option {
	return func(o *Options) { o.Retries = r }
}
func Trans(t *http.Transport) Option {
	return func(o *Options) { o.Trans = t }
}
func Insecure(b bool) Option {
	return func(o *Options) {
		o.Insecure = b
	}
}
func TLSConfig(t *tls.Config) Option {
	return func(o *Options) { o.TLSConfig = t }
}

type CallOptions struct {
	Success     SuccessCallback     // 请求成功回调
	NotModified NotModifiedCallback // 请求未修改回调
	Before      RequestBefore       // 请求之前处理request
	Reties      int32               // 重试次数
	Timeout     time.Duration
}

func newCallOptions(opt ...CallOption) CallOptions {
	opts := CallOptions{}
	for _, o := range opt {
		o(&opts)
	}
	return opts
}

func WithSuccess(s SuccessCallback) CallOption {
	return func(o *CallOptions) { o.Success = s }
}
func WithNotModified(n NotModifiedCallback) CallOption {
	return func(o *CallOptions) { o.NotModified = n }
}
func WithBefore(r RequestBefore) CallOption {
	return func(o *CallOptions) { o.Before = r }
}
func WithRetries(r int32) CallOption {
	return func(o *CallOptions) { o.Reties = r }
}
func WithTimeout(t time.Duration) CallOption {
	return func(o *CallOptions) { o.Timeout = t }
}
