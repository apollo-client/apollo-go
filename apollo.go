package apollo

import (
	"errors"
	"unsafe"
)

type Config struct {
	AppId         string `json:"appId"`
	Cluster       string `json:"cluster"`
	NamespaceName string `json:"namespaceName"`
	Secret        string `json:"secret"`
	Addr          string `json:"addr"`
}

var (
	conf *Config
)

func Init(c *Config) error {
	if c == nil {
		return errors.New("config nil")
	}
	conf = c
	go StartNotifications()
	return nil
}

func Watch(namespace string, deft interface{}, ptr unsafe.Pointer) error {
	return nil
}
