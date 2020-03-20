package yaml

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
)

//Yaml2Json Convert yaml to json data
func Yaml2Json(yamldata []byte) (jsondata []byte, err error) {
	m := map[interface{}]interface{}{}
	err = yaml.Unmarshal(yamldata, &m)
	if err != nil {
		return nil, err
	}
	jsonStruct, err := jsonConvert(m)
	if err != nil {
		return nil, err
	}
	return json.Marshal(jsonStruct)
}

// encoding/json not support map[interface{}]interface{} need a convert
func jsonConvert(m interface{}) (interface{}, error) {
	switch v := m.(type) {
	case map[interface{}]interface{}:
		res := map[string]interface{}{}
		for k, v2 := range v {
			convertv, err := jsonConvert(v2)
			if err != nil {
				return nil, err
			}
			switch k2 := k.(type) {
			case string:
				res[k2] = convertv
			default:
				return nil, fmt.Errorf("unsupport map key type:%T", k)
			}
		}
		return res, nil
	case []interface{}:
		res := make([]interface{}, len(v))
		for i, v2 := range v {
			convertv, err := jsonConvert(v2)
			if err != nil {
				return nil, err
			}
			res[i] = convertv
		}
		return res, nil
	default:
		return m, nil
	}
}
