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

func Request(reqURL string) (int, []byte, error) {
	retry := 0

	for {
		retry++
		if retry > 5 {
			break
		}
		status, body, err := doRequest(httpClient, reqURL)
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

func doRequest(c *http.Client, reqURL string) (int, []byte, error) {
	resp, err := c.Get(reqURL)
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
