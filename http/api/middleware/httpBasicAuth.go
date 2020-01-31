package middleware

import (
	"encoding/base64"
	//"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xieqiaoyu/xin/http/api"
)

//HTTPBasicAuthVerifyFunc 传入函数的第一个参数是用户名，第二个参数名是密码, 第二个返回参数是用户,会被设置在全局变量中
type HTTPBasicAuthVerifyFunc func(string, string) (pass bool, user interface{})

//HTTPBasicAuth 参考 rfc7235和rfc2617 实现接口http Basic 认证的中间键
// 需要传入一个认证函数来实现函数的认证
// 这个中间键应该放在渲染中间键之后
func HTTPBasicAuth(verifyFunc HTTPBasicAuthVerifyFunc) gin.HandlerFunc {
	//alert := "Authorization Required"
	//realm := "Basic realm=" + strconv.Quote(alert)
	// 解析header 的正则对象
	return func(c *gin.Context) {
		var userName, userPass string
		authString := c.GetHeader("Authorization")
		if len(authString) > 6 && strings.ToUpper(authString[0:6]) == "BASIC " {
			//TODO: 过滤多余的空格
			authorization, err := base64.StdEncoding.DecodeString(authString[6:])
			if err == nil {
				authorizationArray := strings.Split(string(authorization), ":")
				userName = authorizationArray[0]
				userPass = authorizationArray[1]
			}
		}
		pass, user := verifyFunc(userName, userPass)
		if !pass {
			//当返回这个值的时候浏览器会默认弹出一个密码框，非常影响体验，暂时先不返回
			//c.Header("WWW-Authenticate", realm)
			// 设置http 401 错误
			c.Set(api.StatusKey, 401)
			c.Set(api.ErrKey, "Auth Fail")
			c.Abort()
			return
		}
		c.Set(api.UserKey, user)
		c.Next()
	}
}
