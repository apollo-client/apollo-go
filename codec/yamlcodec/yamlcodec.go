package yamlcodec

import (
	"encoding/json"
	"reflect"

	"github.com/xnzone/apollo-go/codec"
	"gopkg.in/yaml.v2"
)

type yamlCodec struct{}

func (c *yamlCodec) Parse(configurations map[string]json.RawMessage, deft map[string]json.RawMessage, _ reflect.Type) (map[string]interface{}, error) {
	var str string
	_ = json.Unmarshal(configurations["content"], &str)
	res := make(map[string]interface{})
	_ = yaml.Unmarshal([]byte(str), &res)
	for k, v := range deft {
		if _, ok := res[k]; !ok {
			res[k] = v
		}
	}
	return res, nil
}

func NewCodec() codec.Codec {
	return &yamlCodec{}
}
