//Package jsonschema 提供利用 jsonschema 来进行json 验证的一些工具函数
package jsonschema

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/xeipuuv/gojsonschema"
)

var (
	cache = new(sync.Map)
)

func stringToSchema(str string) (*gojsonschema.Schema, error) {
	// 用schema 字符串的sha1 散列作为缓存的索引
	// 必须先把值赋给一个变量才能将[]byte 转换为字符串
	schemaHashBytes := sha1.Sum([]byte(str))
	schema, cached := cache.Load(schemaHashBytes)
	if !cached {
		schemaTmp, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(str))
		if err != nil {
			return nil, err
		}
		schema, _ = cache.LoadOrStore(schemaHashBytes, schemaTmp)
		// 这里直接return 防止后续改动造成scheme 闭包刷新的问题
		return schema.(*gojsonschema.Schema), nil
	}
	return schema.(*gojsonschema.Schema), nil
}

//ValidJSONString 给定JSON 字符串和 json schema 字符串,验证 json 字符串合法性
// 第一个返回参数是json 是否合法第二个参数返回有无错误发生
func ValidJSONString(jsonStr, schemaStr string) (bool, error) {
	schema, err := stringToSchema(schemaStr)
	if err != nil {
		// schema 不应该出现问题
		panic(fmt.Sprintf("Build schema Fail,err:%s schema:%v", err, schemaStr))
	}
	result, err := schema.Validate(gojsonschema.NewStringLoader(jsonStr))
	// 传入的json 串是错的有可能
	if err != nil {
		return false, err
	}
	if !result.Valid() {
		var ers []string
		for _, e := range result.Errors() {
			ers = append(ers, e.String())
		}
		return false, errors.New(strings.Join(ers, "; "))
	}
	return true, nil
}
