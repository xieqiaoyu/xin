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

//VerifyAndUnmarshalReqestBodyAsJSON 验证gin context 中的request body 是否是合法的json ,如果合法直接将其解析到传入的结构体中
// 解析发生的错误会返回在 error 中
func verifyAndUnmarshalReqestBodyAsJSON(c *gin.Context, schemaStr string, obj interface{}) error {
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
	err := verifyAndUnmarshalReqestBodyAsJSON(c, schemaStr, obj)
	if err != nil {
		//TODO:目前认为context 直接在这个地方abort 不太好，这样函数就被限制得太死了
		SetErrorf("Check JSON err: %w", err).Apply(c)
		return false, err
	}
	return true, nil
}

//SetStatus get an Cache with giving response body code
func SetStatus(code int) *Cache {
	return &Cache{
		Status: &code,
	}
}

//SetData get an Cache with giving data
func SetData(data interface{}) *Cache {
	return &Cache{
		Data: data,
	}
}

//SetError get an Cache with giving error
func SetError(err error) *Cache {
	return &Cache{
		Err: err,
	}
}

//SetErrorf get an Cache with error by giving format
func SetErrorf(format string, a ...interface{}) *Cache {
	var err error
	if len(a) > 0 {
		err = fmt.Errorf(format, a...)
	} else {
		err = errors.New(format)
	}
	return SetError(err)
}

//Cache  response cache of an API Call
type Cache struct {
	Status *int
	Data   interface{}
	Err    error
}

//Apply apply apiContext on given gin context ,this make the real change of an api call
func (ac *Cache) Apply(c *gin.Context) {
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

//SetStatus set the response code of the Cache
func (ac *Cache) SetStatus(code int) *Cache {
	ac.Status = &code
	return ac
}

//SetData set the response data of the Cache
func (ac *Cache) SetData(data interface{}) *Cache {
	ac.Data = data
	return ac
}

//SetError set the error of the Cache
func (ac *Cache) SetError(err error) *Cache {
	ac.Err = err
	return ac
}

//SetErrorf set the error of the Cache by given string format
func (ac *Cache) SetErrorf(format string, a ...interface{}) *Cache {
	var err error
	if len(a) > 0 {
		err = fmt.Errorf(format, a...)
	} else {
		err = errors.New(format)
	}
	return ac.SetError(err)
}
