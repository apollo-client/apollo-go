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

// ApolloAuth apollo default auth
type ApolloAuth struct{}

func (a *ApolloAuth) Headers(uri string, appID string, secret string) map[string][]string {
	headers := make(map[string][]string, 2)
	ms := time.Now().UnixNano() / int64(time.Millisecond)
	ts := strconv.FormatInt(ms, 10)
	query := parse(uri)

	str := fmt.Sprintf("%s\n%s", ts, query)

	sign := sign(str, secret)

	signatures := make([]string, 0, 1)
	signatures = append(signatures, fmt.Sprintf("Apollo %s:%s", appID, sign))
	headers["Authorization"] = signatures

	timestamps := make([]string, 0, 1)
	timestamps = append(timestamps, ts)
	headers["Timestamp"] = timestamps

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
