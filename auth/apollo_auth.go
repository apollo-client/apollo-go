package auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type ApolloAuth struct{}

func NewApolloAuth() *ApolloAuth {
	return &ApolloAuth{}
}

func (a ApolloAuth) Header(rawURL string, appId string, secret string) map[string]string {
	ms := time.Now().UnixNano() / int64(time.Millisecond)
	ts := strconv.FormatInt(ms, 10)

	query := parse(rawURL)
	str := fmt.Sprintf("%s\n%s", ts, query)
	signature := sign(str, secret)

	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Apollo %s:%s", appId, signature)
	headers["Timestamp"] = ts
	return headers
}

func sign(str string, secret string) string {
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(str))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func parse(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	query := u.Path
	if len(u.RawQuery) > 0 {
		query = fmt.Sprintf("%s?%s", query, u.RawQuery)
	}
	return query
}
