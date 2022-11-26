package apollo

import (
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	"github.com/xnzone/apollo-go/transport"
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
		Addr:       "http://81.68.181.139:8080",
		AppId:      "apollo-go",
		Secret:     "",
		Cluster:    "DEV",
		IsBackup:   true,
		BackupPath: "./",
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
	c, _ := NewClient(apolloApp, Transport(trans))
	err := c.Watch("application", mDeft, &mPtr)
	t.Log(err)
	t.Log(mDeft)
	err = c.Watch("app-common", mComDeft, &mComPtr)
	t.Log(err)
	t.Log(mComDeft)
	for i := 0; i < 10; i++ {
		t.Log(DC())
		t.Log(DCom())
		time.Sleep(1 * time.Second)
	}
}
