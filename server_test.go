package apollo

import "testing"

var (
	testApp = &Application{
		Addr:    "http://81.68.181.139:8080",
		AppId:   "apollo-go",
		Secret:  "",
		Cluster: "DEV",
	}
)

func TestGetConfigs(t *testing.T) {
	status, body, err := GetConfigs(testApp, "application", "")
	t.Log(status, body, err)
}
