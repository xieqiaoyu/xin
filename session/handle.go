package session

//Handle session 的处理函数
type Handle interface {
	Load(sessionID string) (session Session, found bool, err error)
	Save(sessionID string, session Session) (ttl int, err error)
}
