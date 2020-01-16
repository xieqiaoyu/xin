package session

//Handler session Handler
type Handler interface {
	Load(sessionID string) (session Session, found bool, err error)
	Save(sessionID string, session Session) (ttl int, err error)
}
