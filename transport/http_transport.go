package transport

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"
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

func (h *HTTPTransport) Do(reqURL string, opt ...Option) (int, []byte, error) {
	h.Init(opt...)
	if h.opts.Timeout > 0 {
		h.opts.Client.Timeout = h.opts.Timeout
	}
	if h.opts.Trans != nil {
		h.opts.Client.Transport = h.opts.Trans
	}
	retry := 0

	for {
		retry++
		if retry > h.opts.MaxRetries {
			break
		}
		status, body, err := doRequest(reqURL, h.opts)
		if err != nil {
			time.Sleep(h.opts.RetryInterval)
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

func doRequest(rawURL string, opts *Options) (int, []byte, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return 0, nil, err
	}
	if len(opts.Headers) > 0 {
		for k, v := range opts.Headers {
			req.Header.Set(k, v)
		}
	}
	if opts.Hook != nil {
		if err = opts.Hook(req); err != nil {
			return 0, nil, err
		}
	}
	resp, err := opts.Client.Do(req)
	if err != nil {
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
	if err == io.EOF {
		err = nil
	}
	return resp.StatusCode, body, err
}
