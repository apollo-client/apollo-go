package apollo

import "testing"

func TestRequest(t *testing.T) {
	reqURL := "https://www.baidu.com"
	status, _, err := Request(reqURL)
	t.Log(status, err)
}
