package apollo

import (
	"sync/atomic"
	"testing"
	"time"
	"unsafe"
)

type TestConfig struct {
	Map     map[string]string `json:"map"`
	Struct  Person            `json:"struct"`
	Strings []string          `json:"strings"`
	Ints    []int32           `json:"ints"`
	String  string            `json:"string"`
	Int     int64             `json:"int"`
	Float   float32           `json:"float"`
}

type Person struct {
	Name string `json:"name"`
	Age  int64  `json:"age"`
}

var (
	apolloApp = &Application{
		Addr:    "http://81.68.181.139:8080",
		AppId:   "apollo-go",
		Secret:  "",
		Cluster: "DEV",
	}
	mPtr  unsafe.Pointer
	mDeft = &TestConfig{
		Float: 0.12,
		Ints:  []int32{2, 3},
	}
)

func DC() *TestConfig {
	p := atomic.LoadPointer(&mPtr)
	if nil != p {
		return (*TestConfig)(p)
	}
	return mDeft
}

func TestWatch(t *testing.T) {
	Init(apolloApp)
	err := Watch("application", mDeft, &mPtr)
	t.Log(err)
	t.Log(mDeft)
	for i := 0; i < 1; i++ {
		t.Log(DC())
		time.Sleep(3 * time.Second)
	}
}
