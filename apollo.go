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

	jsoniter "github.com/json-iterator/go"
)

type Application struct {
	AppId   string `json:"appId"`
	Cluster string `json:"cluster"`
	Secret  string `json:"secret"`
	Addr    string `json:"addr"`
}

var (
	app *Application
)

// Init init application config
func Init(c *Application) error {
	if c == nil {
		return errors.New("config nil")
	}
	app = c
	return nil
}

// Watch watch namespace struct
func Watch(namespace string, deft interface{}, ptr *unsafe.Pointer) error {
	cb, err := namespaceCallback(deft, ptr)
	if err != nil {
		return err
	}
	return WatchNamespace(namespace, cb)
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

	return func(ctx context.Context, apol *Apollo) (err error) {
		// apol or apol configurations nil, return
		if apol == nil || apol.Configurations == nil {
			return
		}
		// fill in default value
		tmp := apol.Configurations
		for k, v := range mdeft {
			if _, ok := tmp[k]; !ok {
				tmp[k] = v
			}
		}

		nd := reflect.New(dt.Type())
		nt := reflect.TypeOf(deft).Elem()
		nm := make(map[string]interface{})
		// configurations are string, so use reflect, then marshal and unmarshal
		for num := 0; num < nt.NumField(); num++ {
			key := nt.Field(num).Tag.Get("json")
			typ := nt.Field(num).Type.Kind()
			switch typ {
			case reflect.Struct, reflect.Map, reflect.Array, reflect.Slice:
				var str string
				_ = json.Unmarshal(tmp[key], &str)
				var val reflect.Value
				if typ == reflect.Struct {
					val = reflect.New(nt.Field(num).Type)
				}
				// map and slice not available, do not know why
				if typ == reflect.Map {
					val = reflect.MakeMap(nt.Field(num).Type)
				}
				if typ == reflect.Array || typ == reflect.Slice {
					val = reflect.MakeSlice(nt.Field(num).Type, 0, 10)
				}
				vpt := val.Interface()
				_ = jsoniter.Unmarshal([]byte(str), &vpt)
				nm[key] = val.Interface()
			case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int8, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				var str string
				_ = json.Unmarshal(tmp[key], &str)
				val, _ := strconv.ParseInt(str, 10, 64)
				nm[key] = val
			case reflect.Float32, reflect.Float64:
				var str string
				_ = json.Unmarshal(tmp[key], &str)
				val, _ := strconv.ParseFloat(str, 64)
				nm[key] = val
			case reflect.Bool:
				var str string
				_ = json.Unmarshal(tmp[key], &str)
				val, _ := strconv.ParseBool(str)
				nm[key] = val
			default:
				var str string
				_ = json.Unmarshal(tmp[key], &str)
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
func WatchNamespace(namespace string, cb WatchCallback) error {
	status, apol, err := GetConfigs(app, namespace, "")
	if err != nil || status != http.StatusOK {
		return fmt.Errorf("watch namespace:%s, err:%v", namespace, err)
	}
	if err = safeCallback(&apol, cb); err != nil {
		return fmt.Errorf("watch namespace:%s, err:%v", namespace, err)
	}

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			ns, na, ne := GetConfigs(app, namespace, apol.ReleaseKey)
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
