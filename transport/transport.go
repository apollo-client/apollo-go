package transport

import "net/http"

type SuccessCallback func([]byte) error
type NotModifiedCallback func([]byte) error
type RequestBefore func(req *http.Request) error

type Option func(o *Options)
type CallOption func(o *CallOptions)

type Transport interface {
	Init(opts ...Option) error
	Options() Options
	Do(reqURL string, opts ...CallOption) error
}
