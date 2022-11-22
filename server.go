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
	return fmt.Sprintf("http://%s/configs/%s/%s/%s?releaseKey=%s&ip=%s",
		app.Addr,
		url.QueryEscape(app.AppId),
		url.QueryEscape(app.Cluster),
		url.QueryEscape(namespace),
		url.QueryEscape(releaseKey),
		GetLocalAddr())
}

func GetConfigs(app *Application, namespace string, releaseKey string) (status int, apol *Apollo, err error) {
	reqURL := ConfigsURL(app, namespace, releaseKey)
	var body []byte
	status, body, err = Request(reqURL)
	if err != nil {
		return
	}
	if status != http.StatusOK {
		return
	}
	err = json.Unmarshal(body, apol)
	if err != nil {
		return
	}
	return
}
