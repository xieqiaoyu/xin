package session

import (
	"github.com/xieqiaoyu/xin"
	"sync"
)

//Memory memor Session implement for test perpose
type Memory struct {
	innerMap sync.Map
}

//Load Session handle implement
func (s *Memory) Load(sessionID string) (Session, bool, error) {
	v, ok := s.innerMap.Load(sessionID)
	if !ok {
		return nil, false, nil
	}
	session, ok := v.(Session)
	if !ok {
		return nil, true, xin.WrapEf(&xin.InternalError{}, "session id %s  is not a valid session interface", sessionID)
	}
	return session, true, nil
}

//Save Session handle implement
func (s *Memory) Save(sessionID string, session Session) (ttl int, err error) {
	s.innerMap.Store(sessionID, session)
	return 0, nil
}
