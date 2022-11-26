package apollo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"sync/atomic"
	"time"
	"unsafe"
)

type Application struct {
	AppId      string `json:"appId"`
	Cluster    string `json:"cluster"`
	Secret     string `json:"secret"`
	Addr       string `json:"addr"`
	IsBackup   bool   `json:"is_backup"`
	BackupPath string `json:"backup_path"`
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
	cb, err := namespaceCallback(deft, ptr)
	if err != nil {
		return err
	}
	return c.watchNamespace(namespace, cb)
}

// namespaceCallback namespace callback function
func namespaceCallback(deft interface{}, ptr *unsafe.Pointer) (WatchCallback, error) {
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
		tmp := apol.Configurations
		def := make(map[string]bool)
		for k, v := range mdeft {
			if _, ok := tmp[k]; !ok {
				tmp[k] = v
				def[k] = true
			}
		}

		nd := reflect.New(dt.Type())
		nt := reflect.TypeOf(deft).Elem()
		nm := make(map[string]interface{})
		// configurations are string, so use reflect, then marshal and unmarshal
		for num := 0; num < nt.NumField(); num++ {
			key := nt.Field(num).Tag.Get("json")
			typ := nt.Field(num).Type.Kind()
			// json.RawMessage to string
			var str string
			if _, ok := def[key]; ok {
				str = string(tmp[key])
			} else {
				_ = json.Unmarshal(tmp[key], &str)
			}
			switch typ {
			case reflect.Struct:
				val := reflect.New(nt.Field(num).Type)
				vpt := val.Interface()
				_ = json.Unmarshal([]byte(str), &vpt)
				nm[key] = val.Interface()
			case reflect.Array, reflect.Slice:
				var val []interface{}
				_ = json.Unmarshal([]byte(str), &val)
				nm[key] = val
			case reflect.Map:
				val := make(map[string]interface{})
				_ = json.Unmarshal([]byte(str), &val)
				nm[key] = val
			case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int8, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				val, _ := strconv.ParseInt(str, 10, 64)
				nm[key] = val
			case reflect.Float32, reflect.Float64:
				val, _ := strconv.ParseFloat(str, 64)
				nm[key] = val
			case reflect.Bool:
				val, _ := strconv.ParseBool(str)
				nm[key] = val
			default:
				nm[key] = string(str)
			}
		}
		// marshal and unmarshal
		tbs, _ := json.Marshal(nm)
		_ = json.Unmarshal(tbs, nd.Interface())

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
