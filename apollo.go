package apollo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/xnzone/apollo-go/codec"
	"github.com/xnzone/apollo-go/codec/jsoncodec"
	"github.com/xnzone/apollo-go/codec/properties"
	"github.com/xnzone/apollo-go/codec/yamlcodec"
)

type Application struct {
	AppId   string `json:"appId"`
	Cluster string `json:"cluster"`
	Secret  string `json:"secret"`
	Addr    string `json:"addr"`
}

// Client apollo client
type Client struct {
	App  *Application // application config
	opts *Options     // options
}

// NewClient new apollo client
func NewClient(c *Application, opt ...Option) (*Client, error) {
	if c == nil {
		return nil, errors.New("config nil")
	}
	opts := newOptions(opt...)
	return &Client{
		App:  c,
		opts: opts,
	}, nil
}

// Watch watch namespace struct
func (c *Client) Watch(namespace string, deft interface{}, ptr *unsafe.Pointer) error {
	var code codec.Codec
	ext := namespace[strings.LastIndex(namespace, ".")+1:]
	switch ext {
	case "json":
		code = jsoncodec.NewCodec()
	case "yaml", "yml":
		code = yamlcodec.NewCodec()
	case "xml":
		return errors.New("not support xml namespace")
	case "txt":
		return errors.New("not support txt namespace")
	default:
		code = properties.NewCodec()
	}
	cb, err := namespaceCallback(deft, ptr, code)
	if err != nil {
		return err
	}
	return c.watchNamespace(namespace, cb)
}

// namespaceCallback namespace callback function
func namespaceCallback(deft interface{}, ptr *unsafe.Pointer, code codec.Codec) (WatchCallback, error) {
	if reflect.Ptr != reflect.TypeOf(deft).Kind() {
		return nil, errors.New("default must be a pointer")
	}
	// store default pointer
	dt := reflect.ValueOf(deft).Elem()
	ua := unsafe.Pointer(dt.UnsafeAddr())
	atomic.StorePointer(ptr, ua)

	// default value map
	var mdeft map[string]json.RawMessage
	bs, _ := json.Marshal(deft)
	if err := json.Unmarshal(bs, &mdeft); err != nil {
		return nil, err
	}

	return func(_ context.Context, apol *Apollo) (err error) {
		// apol or apol configurations nil, return
		if apol == nil || apol.Configurations == nil {
			return
		}
		// fill in default value

		nd := reflect.New(dt.Type())
		nt := reflect.TypeOf(deft).Elem()
		nm, err := code.Parse(apol.Configurations, mdeft, nt)
		if err != nil {
			return err
		}
		// marshal and unmarshal
		tbs, _ := json.Marshal(nm)
		err = json.Unmarshal(tbs, nd.Interface())

		// store new pointer
		nptr := unsafe.Pointer(nd.Elem().UnsafeAddr())
		atomic.StorePointer(ptr, nptr)
		return
	}, nil
}

// WatchCallback watch callback define
type WatchCallback func(ctx context.Context, apol *Apollo) error

// WatchNamespace watch namespace and callback
func (c *Client) watchNamespace(namespace string, cb WatchCallback) error {
	status, apol, err := c.getConfigs(namespace, "")
	if err != nil || status != http.StatusOK {
		return fmt.Errorf("watch namespace:%s, err:%v", namespace, err)
	}
	if err = safeCallback(&apol, cb); err != nil {
		return fmt.Errorf("watch namespace:%s, err:%v", namespace, err)
	}

	go func() {
		ticker := time.NewTicker(c.opts.WatchInterval)
		for range ticker.C {
			ns, na, ne := c.getConfigs(namespace, apol.ReleaseKey)
			if ne != nil || ns != http.StatusOK {
				continue
			}
			apol = na
			_ = safeCallback(&apol, cb)
		}
	}()
	return nil
}

// safeCallback recover if callback failed
func safeCallback(apol *Apollo, cb WatchCallback) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var msg string
			switch e := r.(type) {
			case string:
				msg = e
			case error:
				msg = err.Error()
			default:
				msg = "unknown panic type"
			}
			err = errors.New("callback panic:" + msg)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = cb(ctx, apol); err != nil {
		return fmt.Errorf("callback failed err:%v", err)
	}
	return
}
