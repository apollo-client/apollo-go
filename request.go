package apollo

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	httpClient = &http.Client{}
)

type Option func(o *Options)

type Options struct {
	Headers map[string]string
}

func newOptions(opts ...Option) Options {
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
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

func Request(reqURL string, opts ...Option) (int, []byte, error) {
	opt := newOptions(opts...)
	retry := 0

	for {
		retry++
		if retry > 5 {
			break
		}
		status, body, err := doRequest(httpClient, reqURL, &opt)
		if err != nil {
			continue
		}
		if status != http.StatusOK && status != http.StatusNotModified {
			continue
		}
		return status, body, err
	}
	var err error
	if retry > 5 {
		err = errors.New("over max retry still error")
	}
	return 0, nil, err
}

func doRequest(c *http.Client, reqURL string, opt *Options) (int, []byte, error) {
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return 0, nil, err
	}
	if len(opt.Headers) > 0 {
		for k, v := range opt.Headers {
			req.Header.Set(k, v)
		}
	}
	resp, err := c.Do(req)
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
