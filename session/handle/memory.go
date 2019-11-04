package handle

import (
	"github.com/xieqiaoyu/xin"
	xsession "github.com/xieqiaoyu/xin/session"
	"sync"
)

//MemorySession 内存session 用于测试
type Memory struct {
	innerMap sync.Map
}

//Load Load
func (s *Memory) Load(sessionID string) (xsession.Session, bool, error) {
	v, ok := s.innerMap.Load(sessionID)
	if !ok {
		return nil, false, nil
	}
	session, ok := v.(xsession.Session)
	if !ok {
		return nil, true, xin.WrapEf(&xin.InternalError{}, "session id %s  is not a valid session interface", sessionID)
	}
	return session, true, nil
}

//Save Save
func (s *Memory) Save(sessionID string, session xsession.Session) (ttl int, err error) {
	s.innerMap.Store(sessionID, session)
	return 0, nil
}
