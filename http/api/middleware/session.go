package middleware

import (
	"encoding/base32"
	"github.com/gin-gonic/gin"
	"github.com/xieqiaoyu/xin"
	"github.com/xieqiaoyu/xin/http/api"
	xlog "github.com/xieqiaoyu/xin/log"
	xsession "github.com/xieqiaoyu/xin/session"
	"github.com/xieqiaoyu/xin/util/random"
	"strings"
)

func generateSessionID() string {
	return strings.TrimRight(
		base32.StdEncoding.EncodeToString(
			random.Bytes(32)), "=")
}

//Session middleware handle session issue
func Session(name string, handler xsession.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, _ := c.Cookie(name)
		if sessionID != "" {
			session, exists, err := handler.Load(sessionID)
			if err == nil && exists {
				c.Set(api.SessionKey, session)
			} else {
				// 出错和未找到的情况均将session id 置空
				sessionID = ""
				if err != nil {
					xlog.Errorf("Load session err:%s", err)
				}
			}
		}
		c.Next()

		session, exists := c.Get(api.SessionKey)
		if exists {
			if session != nil {
				newSession, ok := session.(xsession.Session)
				if ok {
					if sessionID == "" {
						sessionID = generateSessionID()
					}
					ttl, err := handler.Save(sessionID, newSession)
					if err != nil {
						xlog.Errorf("Save session err:%s", err)
					} else {
						host := c.GetHeader("Host")
						c.SetCookie(name, sessionID, ttl, "", host, false, true)
					}
				} else {
					xlog.Warningf("session in context type is %v not a session interface, url:%s", session, c.Request.URL)
				}
			}
		}
	}
}

//SessionUserAuth Get a Session User auth middleware , this should after Session middleware
// userKey is the key save user data
func SessionUserAuth(userKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionStruct, exists := c.Get(api.SessionKey)
		if exists && sessionStruct != nil {
			session, ok := sessionStruct.(xsession.Session)
			if !ok {
				c.Set(api.StatusKey, 500)
				c.Set(api.ErrKey, xin.NewWrapEf("malformed session %T is not a Session interface", sessionStruct))
				c.Abort()
				return
			}
			// 按login
			user, exists := session.Get(userKey)
			if exists && user != nil {
				c.Set(api.UserKey, user)
				return
			}
		}
		c.Set(api.StatusKey, 401)
		c.Set(api.ErrKey, "Unauthorize Access")
		c.Abort()
		return
	}
}
