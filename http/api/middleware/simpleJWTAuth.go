package middleware

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	jwt_request "github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"github.com/xieqiaoyu/xin/http/api"
)

//JWTAuthConfig JWT授权验证的配置结构体
type JWTAuthConfig struct {
	Keyfunc jwt.Keyfunc
	Options []jwt_request.ParseFromRequestOption
}

//JWTAuthConfigFunc 提供JWT 验证配置的函数，这个主要为了保证验证的线程安全
type JWTAuthConfigFunc func() *JWTAuthConfig

//SimpleJWTAuth 实现接口基于JWT 的 http Bearer 认证的中间键,参考 rfc6750 和rfc7519
// 如果认证成功 会在context Set token 为后续操作服务
func SimpleJWTAuth(configFunc JWTAuthConfigFunc) gin.HandlerFunc {
	realm := "Bearer realm=%s,error=\"invalid_token\",error_description=\"%s\""
	return func(c *gin.Context) {
		var errMsg string
		config := configFunc()
		token, err := jwt_request.ParseFromRequest(c.Request, jwt_request.AuthorizationHeaderExtractor, config.Keyfunc, config.Options...)
		if err == nil && token.Valid {
			// 设置handle 供后面的逻辑取用
			c.Set(api.JWTTokenKey, token)
			c.Next()
			return
		}
		errMsg = fmt.Sprintf("JWT token validate fail,%s", err)
		c.Header("WWW-Authenticate", fmt.Sprintf(realm, "API", errMsg))
		// 设置http 401 错误
		c.Set(api.StatusKey, 401)
		c.Set(api.ErrKey, errMsg)
		// 这里不进行渲染
		c.Abort()
		return
	}
}
