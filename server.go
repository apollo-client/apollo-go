package apollo

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Apollo struct {
	AppId          string                 `json:"appId"`
	Cluster        string                 `json:"cluster"`
	NamespaceName  string                 `json:"namespaceName"`
	Configurations map[string]interface{} `json:"configurations"`
	ReleaseKey     string                 `json:"releaseKey"`
}

type Notification struct {
	NotificationId int64  `json:"notificationId"`
	NamespaceName  string `json:"namespaceName"`
}

func ConfigsURL(conf *Config, releaseKey string) string {
	return fmt.Sprintf("http://%s/configs/%s/%s/%s?releaseKey=%s&ip=%s",
		conf.Addr,
		url.QueryEscape(conf.AppId),
		url.QueryEscape(conf.Cluster),
		url.QueryEscape(conf.NamespaceName),
		url.QueryEscape(releaseKey),
		GetLocalAddr())
}

func NotifyURL(conf *Config, ns []*Notification) string {
	bs, _ := json.Marshal(ns)
	return fmt.Sprintf("http://%s/notifications/v2?appId=%s&cluster=%s&notifications=%s",
		conf.Addr,
		url.QueryEscape(conf.AppId),
		url.QueryEscape(conf.Cluster),
		url.QueryEscape(string(bs)))
}

func ServicesURL(conf *Config) string {
	return fmt.Sprintf("%sservices/config?appId=%s&ip=%s",
		conf.Addr,
		url.QueryEscape(conf.AppId),
		GetLocalAddr())
}

func GetServices(conf *Config, opts ...RequestOption) error {
	reqURL := ServicesURL(conf)
	return Request(reqURL, opts...)
}

func GetNotify(conf *Config, ns []*Notification, opts ...RequestOption) error {
	reqURL := NotifyURL(conf, ns)
	return Request(reqURL, opts...)
}

func GetConfigs(conf *Config, releaseKey string, opts ...RequestOption) error {
	reqURL := ConfigsURL(conf, releaseKey)
	return Request(reqURL, opts...)
}
