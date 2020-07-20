package middleware

import (
	"github.com/gin-gonic/gin"
	xlog "github.com/xieqiaoyu/xin/log"
	"time"
)

func Logger() gin.HandlerFunc {
	// temporary use gin.Logger logic
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()
		statusColor := ""
		methodColer := ""
		resetColor := ""

		statusCode := c.Writer.Status()
		method := c.Request.Method

		latency := time.Now().Sub(start)
		clientIP := c.ClientIP()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()
		if raw != "" {
			path = path + "?" + raw
		}
		xlog.Infof("%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
			statusColor, statusCode, resetColor,
			latency,
			clientIP,
			methodColer, method, resetColor,
			path,
			errorMessage,
		)
	}
}
