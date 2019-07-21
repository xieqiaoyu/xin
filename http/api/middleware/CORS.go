package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//CORS 赋予接口能被浏览器跨域调用的中间键
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 因为没有特别的逻辑需求，对头部的修改直接放在后续逻辑执行之前处理，方便和渲染中间键一起生效
		origin := c.GetHeader("Origin")
		c.Next()
		if origin != "" {
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Accept, Origin, Authorization,ContentType,Referer,X-HTTP-Method-Override")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		}

	}
}

//OptionsOK 将options 请求返回为成功
func OptionsOK() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusOK)
		}
	}
}
