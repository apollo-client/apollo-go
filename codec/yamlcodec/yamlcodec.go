package yamlcodec

import (
	"encoding/json"
	"reflect"

	"github.com/apollo-client/apollo-go/codec"
	"github.com/apollo-client/apollo-go/log"
	"gopkg.in/yaml.v2"
)

type yamlCodec struct{}

func (c *yamlCodec) Parse(configurations map[string]json.RawMessage, deft map[string]json.RawMessage, _ reflect.Type) (map[string]interface{}, error) {
	var str string
	if err := json.Unmarshal(configurations["content"], &str); err != nil {
		log.Errorf("yaml unmarshal err: %v\n", err)
	}
	res := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(str), &res); err != nil {
		log.Errorf("yaml unmarshal err: %v\n", err)
	}
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
