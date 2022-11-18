package apollo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/xnzone/apollo-go/auth"
	"github.com/xnzone/apollo-go/core"
	"github.com/xnzone/apollo-go/transport"
)

var (
	host          = "http://81.68.181.139:8080"
	configs       = "configfiles/json"
	notifications = "notifications/v2"
	appID         = "apollo-go"
	cluster       = "defalut"
	namespaceName = "application"
	trans         = transport.NewHTTPTransport()
)

func TestApollo(t *testing.T) {
	rawURL := fmt.Sprintf("%s/%s/%s/%s/%s", host, configs, appID, cluster, namespaceName)
	success := func(body []byte) error {
		t.Log(string(body))
		return nil
	}
	opts := []transport.CallOption{
		transport.WithSuccess(success),
		transport.WithBefore(func(req *http.Request) error {
			headers := auth.DefaultAuth.Headers(rawURL, appID, "")
			for k, v := range headers {
				for _, vi := range v {
					req.Header.Add(k, vi)
				}
			}
			return nil
		}),
	}
	err := trans.Do(rawURL, opts...)
	t.Log(err)
}

func TestNotifications(t *testing.T) {
	ns := []core.Notification{
		{NamespaceName: namespaceName, NotificationID: 329},
	}
	bs, _ := json.Marshal(&ns)
	rawURL := fmt.Sprintf("%s/%s?appId=%s&cluster=%s&notifications=%s", host, notifications, appID, cluster, url.QueryEscape(string(bs)))
	success := func(body []byte) error {
		t.Log(string(body))
		return nil
	}
	opts := []transport.CallOption{
		transport.WithSuccess(success),
		transport.WithBefore(func(req *http.Request) error {
			headers := auth.DefaultAuth.Headers(rawURL, appID, "")
			t.Log(headers)
			for k, v := range headers {
				for _, vi := range v {
					req.Header.Add(k, vi)
				}
			}
			return nil
		}),
	}
	err := trans.Do(rawURL, opts...)
	t.Log(err)
}
