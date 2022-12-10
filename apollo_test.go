package apollo

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	"github.com/apollo-client/apollo-go/transport"
)

type TestConfig struct {
	Map     map[string]string `json:"map"`
	Struct  Person            `json:"struct"`
	Strings []string          `json:"strings"`
	Ints    []int32           `json:"ints"`
	String  string            `json:"string"`
	Int     int64             `json:"int"`
	Float   float32           `json:"float"`
	Age     int64             `json:"age"`
	Ages    []int64           `json:"ages"`
}

type CommonConfig struct {
	Name    string   `json:"name"`
	Age     int64    `json:"age"`
	Friends []string `json:"friends"`
	Wants   []int32  `json:"wants"`
}

type Person struct {
	Name string `json:"name"`
	Age  int64  `json:"age"`
}

var (
	apolloApp = &Application{
		Addr:    "http://81.68.181.139:8080",
		AppId:   "apollo-go",
		Secret:  "40edd9317add4489a2eaa104054b01e6",
		Cluster: "DEV",
	}
	mPtr  unsafe.Pointer
	mDeft = &TestConfig{
		Float: 0.12,
		Ints:  []int32{2, 3},
		Age:   100,
		Ages:  []int64{1, 2, 3},
	}
	mComDeft = &CommonConfig{}
	mComPtr  unsafe.Pointer
)

func DC() *TestConfig {
	p := atomic.LoadPointer(&mPtr)
	if nil != p {
		return (*TestConfig)(p)
	}
	return mDeft
}

func DCom() *CommonConfig {
	p := atomic.LoadPointer(&mComPtr)
	if p != nil {
		return (*CommonConfig)(p)
	}
	return mComDeft
}

func TestWatch(t *testing.T) {
	trans := transport.NewHTTPTransport(transport.MaxRetries(1))
	c, _ := NewClient(apolloApp, Transport(trans), EnableBackup(true), BackupPath("./"))
	err := c.Watch("application", mDeft, &mPtr)
	t.Log(err)
	t.Log(mDeft)
	err = c.Watch("app-common", mComDeft, &mComPtr)
	t.Log(err)
	t.Log(mComDeft)
	fmt.Println("NumGoroutine:", runtime.NumGoroutine())
	for i := 0; i < 10; i++ {
		fmt.Println("NumGoroutine:", runtime.NumGoroutine())
		t.Log(DC())
		t.Log(DCom())
		time.Sleep(3 * time.Second)
	}
}

func TestWatchJson(t *testing.T) {
	type Json struct {
		Application string `json:"application"`
		Name        string `json:"name"`
		Age         int64  `json:"age"`
	}
	deft := &Json{
		Name: "json",
		Age:  132,
	}
	var ptr unsafe.Pointer

	DC := func() *Json {
		p := atomic.LoadPointer(&ptr)
		if nil != p {
			return (*Json)(p)
		}
		return deft
	}

	c, _ := NewClient(apolloApp, EnableBackup(true), BackupPath("./"))
	err := c.Watch("app-json.json", deft, &ptr)
	t.Log(err)
	for i := 0; i < 10; i++ {
		t.Log(DC())
		time.Sleep(1 * time.Second)
	}
}

func TestWatchYaml(t *testing.T) {
	type Yaml struct {
		Application string `yaml:"application"`
		Name        string `yaml:"name"`
		Age         int64  `yaml:"age"`
	}
	var ptr unsafe.Pointer
	deft := &Yaml{
		Name: "yaml",
	}
	DC := func() *Yaml {
		p := atomic.LoadPointer(&ptr)
		if nil != p {
			return (*Yaml)(p)
		}
		return deft
	}
	c, _ := NewClient(apolloApp, EnableBackup(true), BackupPath("./"))
	err := c.Watch("app-yaml.yml", deft, &ptr)

	t.Log(err)
	for i := 0; i < 10; i++ {
		t.Log(DC())
		time.Sleep(1 * time.Second)
	}
}

func TestNil(t *testing.T) {
	var m map[string]string
	for k, v := range m {
		fmt.Println(k, v)
	}
}
