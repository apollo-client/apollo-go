package auth

import (
	"testing"

	. "github.com/tevid/gohamcrest"
)

var (
	rawURL = "http://baidu.com/a/b?key=1"
	secret = "6ce3ff7e96a24335a9634fe9abca6d51"
	appID  = "testApplication_yang"
)

func TestAuth(t *testing.T) {
	a := ApolloAuth{}
	headers := a.Headers(rawURL, appID, secret)
	Assert(t, headers, HasMapValue("Authorization"))
	Assert(t, headers, HasMapValue("Timestamp"))
}

func TestSign(t *testing.T) {
	s := sign(rawURL, secret)
	Assert(t, s, Equal("mcS95GXa7CpCjIfrbxgjKr0lRu8="))
}

func TestParse(t *testing.T) {
	query := parse(rawURL)
	Assert(t, query, Equal("/a/b?key=1"))
}
