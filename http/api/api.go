//Package api 提供api 相关的一些定义和工具函数
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xieqiaoyu/xin"
	xjsonschema "github.com/xieqiaoyu/xin/util/jsonschema"

	"github.com/gin-gonic/gin"
)

//ValidateRequestBodyJSON 验证gin context 中的request body 是否是合法的json ,如果合法直接将其解析到传入的结构体中
// 解析发生的错误会返回在 error 中
func ValidateRequestBodyJSON(c *gin.Context, schemaStr string, obj interface{}) error {
	requestBody, err := c.GetRawData()
	// 空的request Body 直接返回
	if len(requestBody) == 0 {
		return errors.New("Empty Request Body")
	}
	if err != nil {
		return xin.WrapEf(&xin.InternalError{}, "Fail to load request data: %w", err)
	}
	// 好像不用太关心 json 的valid情况
	_, err = xjsonschema.ValidJSONString(string(requestBody), schemaStr)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(requestBody, obj); err != nil {
		return err
	}
	return nil
}

//CheckReqJSON 检查gin context 中的request body 是否是正确的json ,如果是，则将其解析到传入的结构体中并返回true
// 此函数需要配合responseRender 渲染中间键才会生效
// 如果json 不合法则会直接设置解析错误的内容,并返回false 和错误对象
func CheckReqJSON(c *gin.Context, schemaStr string, obj interface{}) (bool, error) {
	err := ValidateRequestBodyJSON(c, schemaStr, obj)
	if err != nil {
		//TODO:目前认为context 直接在这个地方abort 不太好，这样函数就被限制得太死了
		SetErrorf("Check JSON err: %w", err).Apply(c)
		return false, err
	}
	return true, nil
}

func SetStatus(code int) *ApiContext {
	return &ApiContext{
		Status: &code,
	}
}

func SetData(data interface{}) *ApiContext {
	return &ApiContext{
		Data: data,
	}
}

func SetError(err error) *ApiContext {
	return &ApiContext{
		Err: err,
	}
}

func SetErrorf(format string, a ...interface{}) *ApiContext {
	var err error
	if len(a) > 0 {
		err = fmt.Errorf(format, a...)
	} else {
		err = errors.New(format)
	}
	return SetError(err)
}

type ApiContext struct {
	Status *int
	Data   interface{}
	Err    error
}

func (ac *ApiContext) Apply(c *gin.Context) {
	if ac.Status != nil {
		c.Set(StatusKey, *ac.Status)
	}
	if ac.Data != nil {
		c.Set(DataKey, ac.Data)
	}
	if ac.Err != nil {
		c.Set(ErrKey, ac.Err)
	}
}

func (ac *ApiContext) SetStatus(code int) *ApiContext {
	ac.Status = &code
	return ac
}

func (ac *ApiContext) SetData(data interface{}) *ApiContext {
	ac.Data = data
	return ac
}

func (ac *ApiContext) SetError(err error) *ApiContext {
	ac.Err = err
	return ac
}

func (ac *ApiContext) SetErrorf(format string, a ...interface{}) *ApiContext {
	var err error
	if len(a) > 0 {
		err = fmt.Errorf(format, a...)
	} else {
		err = errors.New(format)
	}
	return ac.SetError(err)
}
