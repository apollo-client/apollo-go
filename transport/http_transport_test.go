package transport

import "testing"

var (
	htrans = NewHTTPTransport()
)

func TestTransport(t *testing.T) {
	sucess := func(body []byte) error {
		t.Log(string(body))
		return nil
	}
	opts := []CallOption{
		WithSuccess(sucess),
	}
	err := htrans.Do("https://ww.google.com", opts...)
	t.Log(err)
}
