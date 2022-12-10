package properties

import (
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/apollo-client/apollo-go/codec"
)

type propertiesCodec struct{}

func (c *propertiesCodec) Parse(configurations map[string]json.RawMessage, deft map[string]json.RawMessage, dt reflect.Type) (map[string]interface{}, error) {
	def := make(map[string]bool)
	for k, v := range deft {
		if _, ok := configurations[k]; !ok {
			configurations[k] = v
			def[k] = true
		}
	}
	nm := make(map[string]interface{})
	// configurations are string, so use reflect, then marshal and unmarshal
	for num := 0; num < dt.NumField(); num++ {
		key := dt.Field(num).Tag.Get("json")
		typ := dt.Field(num).Type.Kind()
		// json.RawMessage to string
		var str string
		if _, ok := def[key]; ok {
			str = string(configurations[key])
		} else {
			_ = json.Unmarshal(configurations[key], &str)
		}
		switch typ {
		case reflect.Struct:
			val := reflect.New(dt.Field(num).Type)
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
	return nm, nil
}

func NewCodec() codec.Codec {
	return &propertiesCodec{}
}
