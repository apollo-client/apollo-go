package jsoncodec

import (
	"encoding/json"
	"reflect"

	"github.com/apollo-client/apollo-go/codec"
	"github.com/apollo-client/apollo-go/log"
)

type jsonCodec struct{}

func (c *jsonCodec) Parse(configurations map[string]json.RawMessage, deft map[string]json.RawMessage, _ reflect.Type) (map[string]interface{}, error) {
	var str string
	if err := json.Unmarshal(configurations["content"], &str); err != nil {
		log.Errorf("json unmarshal err: %v\n", err)
	}
	res := make(map[string]interface{})
	if err := json.Unmarshal([]byte(str), &res); err != nil {
		log.Errorf("json unmarshal err: %v\n", err)
	}

	for k, v := range deft {
		if _, ok := res[k]; !ok {
			res[k] = v
		}
	}
	return res, nil
}

func NewCodec() codec.Codec {
	return &jsonCodec{}
}
