package transport

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type HTTPTransport struct {
	opts *Options
}

func NewHTTPTransport(opts ...Option) *HTTPTransport {
	opt := newOptions(opts...)
	return &HTTPTransport{
		opts: &opt,
	}
}

func (h *HTTPTransport) Init(opts ...Option) error {
	if h.opts == nil {
		opt := newOptions(opts...)
		h.opts = &opt
	} else {
		initOptions(h.opts, opts...)
	}
	return nil
}

func (h *HTTPTransport) Do(reqURL string, opts ...CallOption) error {
	opt := newCallOptions(opts...)
	if h.opts.Client == nil {
		h.opts.Client = &http.Client{}
	}
	h.opts.Client.Timeout = h.opts.Timeout
	if opt.Timeout > 0 {
		h.opts.Client.Timeout = opt.Timeout
	}

	var retries int32 = h.opts.Retries
	if opt.Reties > 0 {
		retries = opt.Reties
	}
	var retry int32 = 0
	for {
		retry++
		if retry > retries {
			break
		}
		err := h.do(reqURL, &opt)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		// break for no error
		break
	}
	if retry > retries {
		return errors.New("over max retry still error")
	}
	return nil
}

func (h *HTTPTransport) do(reqURL string, opt *CallOptions) error {
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return err
	}
	if opt.Before != nil {
		if err = opt.Before(req); err != nil {
			return err
		}
	}

	res, err := h.opts.Client.Do(req)
	if err != nil {
		return err
	}
	if res == nil {
		return errors.New("res nil")
	}
	defer func() { _ = res.Body.Close() }()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	switch res.StatusCode {
	case http.StatusOK:
		if opt.Success != nil {
			return opt.Success(body)
		}
		return nil
	case http.StatusNotModified:
		if opt.NotModified != nil {
			return opt.NotModified(body)
		}
		return nil
	default:
		return errors.New("status error")
	}
}
