# apollo-go
apollo config go

## 介绍
基于携程apollo的go客户端

## 功能

- [x] 动态实时更新通知 通过`Watch`订阅对应的`namespace`变更，当有变更时，保证推送一次给业务
- [x] 支持默认值 当配置值无法获取时，支持业务配置的默认值
- [x] 支持灰度发布 自动获取业务本地ip，当apollo配置中心配置了灰度，会获取灰度配置
- [x] 支持自定义认证和自定义http请求参数
- [ ] 支持本地文件备份 当配置中心出现问题时，通过备份文件加载

## 使用

- Init 初始化apollo配置
- Watch 监听`namespace`结构体

看下面例子

```go
package main


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
    // 默认值配置，配置获取不到时，会用默认值
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

func main() {
    c, _ := apollo.NewClient(apolloApp)
	err := c.Watch("application", mDeft, &mPtr)
    for i := 0; i < 100; i++ {
        fmt.Printf("dconf:%+v", DC())
        time.Sleep(1 * time.Second)
    }
}
```