package api

// define key name used in context
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
	ErrorStatusDefault = 500
	//UserKey api 保存的用户对象的key
	UserKey = "xinUser"
	//JWTTokenKey context 中保存jwt-go token 对象的key
	JWTTokenKey = "xinJWTToken"
	//SessionKey context 中保存session 的 key
	SessionKey = "xinSession"

	//S shorcut of StatusKey
	S = StatusKey
	//E shortcut of ErrKey
	E = ErrKey
	//D shortcut of DataKey
	D = DataKey
	//U shortcut of UserKey
	U = UserKey
)
