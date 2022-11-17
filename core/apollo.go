package core

// Apollo apollo配置
type Apollo struct {
	AppID         string `json:"appId"`
	Cluster       string `json:"cluster"`
	NamespaceName string `json:"namespaceName"`
	ReleaseKey    string `json:"releaseKey"`
}
