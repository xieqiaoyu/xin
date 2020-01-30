package session

//Handler session Handler
type Handler interface {
	//Load a Session by sessionID
	Load(sessionID string) (session Session, found bool, err error)
	//Save a Session with given sessionID
	Save(sessionID string, session Session) (ttl int, err error)
}
