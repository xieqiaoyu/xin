package middleware

import (
	"crypto/rand"
	"encoding/base32"
	"github.com/gin-gonic/gin"
	"github.com/xieqiaoyu/xin/http/api"
	xlog "github.com/xieqiaoyu/xin/log"
	xsession "github.com/xieqiaoyu/xin/session"
	"io"
	"strings"
)

func GenerateRandomKey(length int) []byte {
	k := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return nil
	}
	return k
}

func generateSessionID() string {
	return strings.TrimRight(
		base32.StdEncoding.EncodeToString(
			GenerateRandomKey(32)), "=")
}

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
					xlog.WriteError("Load session err:%s", err)
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
						xlog.WriteError("Save session err:%s", err)
					} else {
						host := c.GetHeader("Host")
						c.SetCookie(name, sessionID, ttl, "", host, false, true)
					}
				} else {
					xlog.WriteWarning("session in context type is %v not a session interface, url:%s", session, c.Request.URL)
				}
			}
		}
	}
}
