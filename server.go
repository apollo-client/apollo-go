package apollo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Apollo struct {
	AppId         string `json:"appId"`
	Cluster       string `json:"cluster"`
	NamespaceName string `json:"namespaceName"`
	ReleaseKey    string `json:"releaseKey"`

	Configurations map[string]json.RawMessage `json:"configurations"`
}

func ConfigsURL(app *Application, namespace string, releaseKey string) string {
	return fmt.Sprintf("%s/configs/%s/%s/%s?releaseKey=%s&ip=%s",
		app.Addr,
		url.QueryEscape(app.AppId),
		url.QueryEscape(app.Cluster),
		url.QueryEscape(namespace),
		url.QueryEscape(releaseKey),
		GetLocalAddr(),
	)
}

func GetConfigs(app *Application, namespace string, releaseKey string) (int, Apollo, error) {
	var apol Apollo
	reqURL := ConfigsURL(app, namespace, releaseKey)
	opts := []Option{
		Headers(GetAuth(reqURL, app.AppId, app.Secret)),
	}
	var body []byte
	status, body, err := Request(reqURL, opts...)
	if err != nil {
		return status, apol, err
	}
	if status != http.StatusOK {
		return status, apol, err
	}
	err = json.Unmarshal(body, &apol)
	if err != nil {
		return status, apol, err
	}
	return status, apol, nil
}
