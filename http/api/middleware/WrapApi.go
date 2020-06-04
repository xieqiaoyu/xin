//Package middleware 自定义 api 中间键包，提供常用业务场景中会使用到的一些中间键封装
// 某些自定义功能可能会需要几个中间键组合在一起才能生效
package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xieqiaoyu/xin"
	"github.com/xieqiaoyu/xin/http/api"
	xlog "github.com/xieqiaoyu/xin/log"
)

// Wrappers api 包装的逻辑map
type Wrappers map[string]gin.HandlerFunc

//WrapperDefaultKey 默认的wrapper 名称
const WrapperDefaultKey = "default"

// SetDefault 设置默认的包装逻辑
func (w Wrappers) SetDefault(wrapper gin.HandlerFunc) Wrappers {
	w[WrapperDefaultKey] = wrapper
	return w
}

//WrapAPI api及中间键逻辑完成之后进行统一的header 和content 设定的中间键
// 使用这个中间键的时候要注意它的执行位置，最好放在所有接口逻辑之前
func WrapAPI(wrappers Wrappers) gin.HandlerFunc {
	return func(c *gin.Context) {
		formarStatus := c.Writer.Status()
		// 如果之前的状态不是200 就不进行包装了
		if formarStatus != 200 {
			return
		}
		c.Next()
		wrapperIndex := WrapperDefaultKey
		userSetWrapperIndex, exists := c.Get(api.WrapperIndex)
		if exists {
			wrapperIndex = userSetWrapperIndex.(string)
		}
		wrapper, found := wrappers[wrapperIndex]
		if !found {
			panic(fmt.Errorf("Wrapper Index %s is not found", wrapperIndex))
		}
		wrapper(c)
	}
}

// SimpleJSONWrapper get a simple json wrapper middware
func SimpleJSONWrapper() gin.HandlerFunc {
	return func(c *gin.Context) {
		status := 200
		apiStatus, statusExists := c.Get(api.StatusKey)
		if statusExists {
			status = apiStatus.(int)
		}
		dataStruct, dataExists := c.Get(api.DataKey)
		if dataExists {
			c.JSON(status, dataStruct)
			return
		}
		c.String(status, "")
	}
}

//XinRESTfulResponse  restful response struct
type XinRESTfulResponse struct {
	Status int    `json:"status"`
	ErrMsg string `json:"err_msg,omitempty"`
}

//XinRESTfulWrapper get a xin restful wrapper middware
// httpstatus = apiStatus%1000
// For example if api status is set to 1404 then api http status will be set to 404
func XinRESTfulWrapper(env xin.Envirment) gin.HandlerFunc {
	return func(c *gin.Context) {
		baseResponseObj := new(XinRESTfulResponse)
		apiStatus, statusExists := c.Get(api.StatusKey)
		errMsg, errExists := c.Get(api.ErrKey)
		dataStruct, dataExists := c.Get(api.DataKey)

		if statusExists {
			baseResponseObj.Status = apiStatus.(int)
		} else if errExists {
			baseResponseObj.Status = api.ErrorStatusDefault
		} else {
			baseResponseObj.Status = api.StatusDefault
		}
		// 根据返回的apiStatus 获取http 的status
		httpStatus := int(baseResponseObj.Status % 1000)
		//TODO: 需要一个更加合理的方式进行有效性的判断
		if httpStatus < 100 || httpStatus >= 600 {
			xlog.Warningf("malformed api status %d", baseResponseObj.Status)
			httpStatus = 500
		}

		if errExists {
			var errMsgString string
			var isInternalError bool
			switch t := errMsg.(type) {
			case string:
				errMsgString = t
			case error:
				errMsgString = t.Error()
				var internalErr *xin.InternalError
				isInternalError = errors.As(t, &internalErr)
			default:
				xlog.Warningf("Unexpected ErrMsg type %T", t)
			}
			//http 状态码 > 500 在正式环境应该屏蔽错误输出并将错误输入到日志中
			if httpStatus >= 500 || isInternalError {
				xlog.Errorf("%s return status %d with error message:%s", c.Request.URL.Path, httpStatus, errMsgString)
				if env.Mode() != xin.ReleaseMode {
					baseResponseObj.ErrMsg = errMsgString
				}
				//TODO: 更加详细的记录包括请求header 和 body
			}
		}

		if dataExists {
			//将data 数据合并到baseResponseObj 中
			var responseMap map[string]interface{}
			var baseResponseMap map[string]interface{}

			// NOTE(xieqiaoyu) 目前不认为下面的json 解析会出错误，暂时不做错误处理
			dataStructJSON, _ := json.Marshal(dataStruct)
			json.Unmarshal(dataStructJSON, &responseMap)
			baseResponseJSON, _ := json.Marshal(baseResponseObj)
			json.Unmarshal(baseResponseJSON, &baseResponseMap)
			// 合并两个map
			// 这种写法在处理上会快一些
			for key, value := range baseResponseMap {
				// baseResponse 不覆盖应用层的定义
				if _, exists := responseMap[key]; !exists {
					responseMap[key] = value
				}
			}
			c.JSON(httpStatus, responseMap)
		} else {
			c.JSON(httpStatus, baseResponseObj)
		}
	}
}

// NewWrappers 创建一个带有default 的wrappers
func NewWrappers() Wrappers {
	return Wrappers{
		WrapperDefaultKey: SimpleJSONWrapper(),
	}
}

// XinRESTfulWrap create a WrapAPI middleware with XinRESTfulWrapper as default wrapper
func XinRESTfulWrap(env xin.Envirment) gin.HandlerFunc {
	wrappers := NewWrappers()
	wrappers.SetDefault(XinRESTfulWrapper(env))
	return WrapAPI(wrappers)
}
