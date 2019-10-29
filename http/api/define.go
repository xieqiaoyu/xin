package api

import (
	"github.com/gin-gonic/gin"
)

const (
	ResultKey = "xinApiResult"
	//StatusKey context 中使用的api 返回 status key 名称
	StatusKey = "xinApiStatus"
	//ErrKey context 中使用的api 返回 error key  名称
	ErrKey = "xinApiErrorMsg"
	//DataKey context 中使用的api 实际返回业务数据的 key 名称
	DataKey = "xinApiData"
	//StatusDefault api handle 未定义status 时返回的默认成功值
	StatusDefault = 200
	//WrapperIndex 指定api 使用的wrapper 逻辑
	WrapperIndex = "xinWrapperIndex"
	//ErrorStatusDefault api handle 未定义status 时返回的默认失败值
	ErrorStatusDefault = 400
	//UserKey api 保存的用户对象的key
	UserKey = "xinUser"
	//JWTTokenKey context 中保存jwt-go token 对象的key
	JWTTokenKey = "xinJWTToken"
	//SessionKey context 中保存session 的 key
	SessionKey = "xinSession"

	//S 用于外部引用StatusKey 的shorcut
	S = StatusKey
	//E 用于外部引用ErrKey 的shortcut
	E = ErrKey
	//D 用于外部引用DataKey 的shortcut
	D = DataKey
	//U 用于外部引用UserKey 的shortcut
	U = UserKey
)

func SetData(c *gin.Context, data interface{}, code ...int) {
	c.Set(DataKey, data)
	if len(code) > 0 {
		c.Set(StatusKey, code[0])
	}
}

func SetError(c *gin.Context, err error, code ...int) {
	c.Set(ErrKey, err)
	if len(code) > 0 {
		c.Set(StatusKey, code[0])
	}
}
