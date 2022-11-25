package transport

import "testing"

func TestRequest(t *testing.T) {
	reqURL := "https://www.baidu.com"
	status, _, err := DefaultTransport.Do(reqURL)
	t.Log(status, err)
}
