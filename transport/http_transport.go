package transport

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/apollo-client/apollo-go/log"
)

type HTTPTransport struct {
	opts *Options
}

func NewHTTPTransport(opt ...Option) *HTTPTransport {
	opts := newOptions(opt...)
	return &HTTPTransport{
		opts: opts,
	}
}

func (h *HTTPTransport) Init(opt ...Option) error {
	if h.opts == nil {
		h.opts = newOptions(opt...)
	} else {
		initOptions(h.opts, opt...)
	}
	return nil
}

func (h *HTTPTransport) Options() *Options {
	return h.opts
}

func (h *HTTPTransport) Do(reqURL string, opt ...CallOption) (int, []byte, error) {
	retry := 0

	for {
		retry++
		if retry > h.opts.MaxRetries {
			break
		}
		status, body, err := h.doRequest(reqURL, opt...)
		if err != nil {
			time.Sleep(h.opts.RetryInterval)
			log.Errorf("do err: %v\n", err)
			continue
		}
		if status != http.StatusOK && status != http.StatusNotModified {
			time.Sleep(h.opts.RetryInterval)
			continue
		}
		return status, body, err
	}
	var err error
	if retry > h.opts.MaxRetries {
		err = errors.New("over max retry still error")
	}
	return 0, nil, err
}

func (h *HTTPTransport) doRequest(rawURL string, opts ...CallOption) (int, []byte, error) {
	opt := &CallOptions{}
	for _, o := range opts {
		o(opt)
	}
	ts := h.opts.Timeout
	if opt.Timeout > 0 {
		ts = opt.Timeout
	}
	if ts > 0 {
		h.opts.Client.Timeout = ts
	}
	if h.opts.Trans != nil {
		h.opts.Client.Transport = h.opts.Trans
	}
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		log.Errorf("http request err: %v\n", err)
		return 0, nil, err
	}
	// call header
	if len(opt.Headers) > 0 {
		for k, v := range opt.Headers {
			req.Header.Set(k, v)
		}
	}
	// client header
	if len(h.opts.Headers) > 0 {
		for k, v := range h.opts.Headers {
			req.Header.Set(k, v)
		}
	}
	if h.opts.Hook != nil {
		if err = h.opts.Hook(req); err != nil {
			log.Errorf("request hook err: %v\n", err)
			return 0, nil, err
		}
	}
	resp, err := h.opts.Client.Do(req)
	if err != nil {
		log.Errorf("doRequest url: %s err: %v\n", rawURL, err)
		return 0, nil, err
	}
	if resp == nil {
		return 0, nil, errors.New("resp nil")
	}
	if resp.Body != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read request body err: %v\n", err)
	}
	if err == io.EOF {
		err = nil
	}
	return resp.StatusCode, body, err
}
