package cc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

func unmarshalYAMLFile(fpath string) (map[string]interface{}, error) {
	content, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	return unmarshalYAML(content)
}

func unmarshalYAML(b []byte) (map[string]interface{}, error) {
	var data map[interface{}]interface{}
	if err := yaml.Unmarshal(b, &data); err != nil {
		return nil, err
	}
	return unknownMapToStringMap(data), nil
}

func unknownMapToStringMap(data map[interface{}]interface{}) map[string]interface{} {
	kv := make(map[string]interface{}, len(data))
	for k, v := range data {
		kv[fmt.Sprintf("%v", k)] = v
	}
	return kv
}

func unmarshalJSONFile(fpath string) (map[string]interface{}, error) {
	content, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	return unmarshalJSON(content)
}

func unmarshalJSON(b []byte) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func toBool(v interface{}, deflt bool) bool {
	if x, ok := v.(bool); ok {
		return x
	}
	return deflt
}

func toString(v interface{}, deflt string) string {
	if x, ok := v.(string); ok {
		return x
	}
	return deflt
}

func toInt(v interface{}, deflt int) int {
	switch x := v.(type) {
	case int:
		return x
	case int32:
		return int(x)
	case float32:
		return int(x)
	case float64: // for JSON
		return int(x)
	case int16:
		return int(x)
	case int8:
		return int(x)
	}
	return deflt
}

func toInt64(v interface{}, deflt int64) int64 {
	switch x := v.(type) {
	case int64:
		return x
	case int:
		return int64(x)
	case int32:
		return int64(x)
	case float64:
		return int64(x)
	case float32:
		return int64(x)
	case int16:
		return int64(x)
	case int8:
		return int64(x)
	}
	return deflt
}

func toFloat64(v interface{}, deflt float64) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case float32:
		return float64(x)
	case int:
		return float64(x)
	case int32:
		return float64(x)
	case int16:
		return float64(x)
	case int8:
		return float64(x)
	}
	return deflt
}
