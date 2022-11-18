package transport

import "testing"

var (
	htrans = NewHTTPTransport()
	host   = "http://81.68.181.139:8080"
)

func TestTransport(t *testing.T) {
	sucess := func(body []byte) error {
		t.Log(string(body))
		return nil
	}
	opts := []CallOption{
		WithSuccess(sucess),
	}
	err := htrans.Do("https://www.google.com", opts...)
	t.Log(err)
}
