package apollo

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type RequestOptions struct {
	Timeout     time.Duration
	Success     SuccessCallback
	NotModified NotModifiedCallback
}

type SuccessCallback func(body []byte) error
type NotModifiedCallback func(body []byte) error

func newRequestOptions(opt ...RequestOption) RequestOptions {
	opts := RequestOptions{}
	for _, o := range opt {
		o(&opts)
	}
	return opts
}

func Timeout(t time.Duration) RequestOption {
	return func(o *RequestOptions) { o.Timeout = t }
}
func Success(fn SuccessCallback) RequestOption {
	return func(o *RequestOptions) { o.Success = fn }
}
func NotModified(fn NotModifiedCallback) RequestOption {
	return func(o *RequestOptions) { o.NotModified = fn }
}

type RequestOption func(o *RequestOptions)

func Request(reqURL string, opts ...RequestOption) error {
	opt := newRequestOptions(opts...)
	c := &http.Client{}
	if opt.Timeout != 0 {
		c.Timeout = opt.Timeout
	} else {
		c.Timeout = 1 * time.Second
	}
	retry := 0

	for {
		retry++
		if retry > 5 {
			break
		}
		if err := doRequest(c, reqURL, &opt); err != nil {
			continue
		} else {
			return nil
		}

	}
	if retry > 5 {
		return errors.New("over max retry still error")
	}
	return nil
}

func doRequest(c *http.Client, reqURL string, opts *RequestOptions) error {
	resp, err := c.Get(reqURL)
	if err != nil {
		return err
	}
	if resp == nil {
		return errors.New("resp nil")
	}
	if resp.Body != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		if opts.Success != nil {
			return opts.Success(body)
		}
		return nil
	case http.StatusNotModified:
		if opts.NotModified != nil {
			return opts.NotModified(body)
		}
		return nil
	default:
		return errors.New("error response")
	}
}
