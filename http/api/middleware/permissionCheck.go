package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/xieqiaoyu/xin/http/api"
)

//CheckPermitFunc 给定一个 context 和一个api 名称,获取request 来源是否有权限访问api
// 注意，不应该在这个函数中修改 c 的内容, 不会影响主逻辑后和其他中间键的context
type CheckPermitFunc func(c *gin.Context, apiName string) (pass bool, user interface{}, err error)

//PermissionCheck 检查权限的中间键，可以提供多个checkpoints 函数进行权限检查,任意一个检查点未通过都会直接返回403,checkpoint 中的Context是独立的，不会影响实际的request context ，checkpoints 间的context 是共享的
func PermissionCheck(apiName string, checkpoints ...CheckPermitFunc) gin.HandlerFunc {
	// 这个中间键应该放在渲染中间键之前
	return func(c *gin.Context) {
		// 在checkpoints 中传递复制避免一些意外情况
		copyC := c.Copy()
		for _, check := range checkpoints {
			pass, user, err := check(copyC, apiName)
			if !pass {
				// 设置http 403 错误
				c.Set(api.StatusKey, 403)
				c.Set(api.ErrKey, err)
				c.Abort()
				return
			}
			if user != nil {
				c.Set(api.UserKey, user)
				copyC.Set(api.UserKey, user)
			}
		}
		c.Next()

	}
}
