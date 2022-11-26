package codec

import (
	"encoding/json"
	"reflect"
)

type Codec interface {
	Parse(configurations map[string]json.RawMessage, deft map[string]json.RawMessage, dt reflect.Type) (map[string]interface{}, error)
}
