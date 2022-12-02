package apollo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/xnzone/apollo-go/log"
	"github.com/xnzone/apollo-go/transport"
	"github.com/xnzone/apollo-go/util"
)

type Apollo struct {
	AppId         string `json:"appId"`
	Cluster       string `json:"cluster"`
	NamespaceName string `json:"namespaceName"`
	ReleaseKey    string `json:"releaseKey"`

	Configurations map[string]json.RawMessage `json:"configurations"`
}

type Notifcation struct {
	NamespaceName string `json:"namespaceName"`
	NotifcationID int64  `json:"notificationId"`
}

func configsURL(app *Application, namespace string, releaseKey string) string {
	return fmt.Sprintf("%s/configs/%s/%s/%s?releaseKey=%s&ip=%s",
		app.Addr,
		url.QueryEscape(app.AppId),
		url.QueryEscape(app.Cluster),
		url.QueryEscape(namespace),
		url.QueryEscape(releaseKey),
		util.GetLocalAddr(),
	)
}

func (c *Client) getConfigs(namespace string, releaseKey string) (int, Apollo, error) {
	var apol Apollo
	app := c.App
	reqURL := configsURL(app, namespace, releaseKey)
	header := transport.Headers(c.opts.Auth.Header(reqURL, app.AppId, app.Secret))
	var body []byte
	status, body, err := c.opts.Transport.Do(reqURL, header)
	if err != nil {
		log.Errorf("get configs namespace: %s, release:%s, err: %v\n", namespace, releaseKey, err)
	}
	// new body must write to backup
	if c.opts.EnableBackup && err == nil && status == http.StatusOK && body != nil {
		filePath := path.Join(c.opts.BackupPath, fmt.Sprintf("%s-%s-%s", c.App.AppId, c.App.Cluster, namespace))
		_ = c.opts.Backup.Write(filePath, body)
	}
	// request failed, read from backup
	if c.opts.EnableBackup && (err != nil || status != http.StatusOK) {
		filePath := path.Join(c.opts.BackupPath, fmt.Sprintf("%s-%s-%s", c.App.AppId, c.App.Cluster, namespace))
		body, err = c.opts.Backup.Read(filePath)
		// read success, set status ok
		if err == nil && body != nil {
			status = http.StatusOK
		}
	}
	if err != nil {
		log.Errorf("get configs namespace: %s, release:%s, err: %v\n", namespace, releaseKey, err)
		return status, apol, err
	}
	if status != http.StatusOK {
		return status, apol, err
	}
	err = json.Unmarshal(body, &apol)
	if err != nil {
		log.Errorf("get configs namespace: %s, release:%s, err: %v\n", namespace, releaseKey, err)
		return status, apol, err
	}
	return status, apol, nil
}

func notificationURL(app *Application, ns []*Notifcation) string {
	bs, _ := json.Marshal(ns)
	return fmt.Sprintf("%s/notifications/v2?appId=%s&cluster=%s&notifications=%s",
		app.Addr,
		url.QueryEscape(app.AppId),
		url.QueryEscape(app.Cluster),
		url.QueryEscape(string(bs)))
}

func (c *Client) getNotifications(ns []*Notifcation) (int, []*Notifcation, error) {
	app := c.App
	reqURL := notificationURL(app, ns)
	header := c.opts.Auth.Header(reqURL, app.AppId, app.Secret)
	var body []byte
	status, body, err := c.opts.Transport.Do(reqURL, transport.Headers(header), transport.Timeout(10*time.Minute))
	if err != nil {
		log.Errorf("get notifications url: %s, err: %v\n", reqURL, err)
		return status, nil, err
	}
	if status != http.StatusOK {
		return status, nil, err
	}
	res := make([]*Notifcation, 0)
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Errorf("get notifications url: %s, err: %v\n", reqURL, err)
		return status, res, err
	}
	return status, res, nil
}
