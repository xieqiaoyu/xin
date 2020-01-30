package session

import (
	"github.com/mediocregopher/radix/v3"
	"github.com/xieqiaoyu/xin"
)

//StorableSessionGenerater StorableSessionGenerater
type StorableSessionGenerater func() StorableSession

//Redis  redis session handler
type Redis struct {
	TTL              int
	store            radix.Client
	sessionGenerater StorableSessionGenerater
}

//NewRedisHandler NewRedisHandler
func NewRedisHandler(store radix.Client, ttl int, generater StorableSessionGenerater) *Redis {
	return &Redis{
		TTL:              ttl,
		store:            store,
		sessionGenerater: generater,
	}
}

func xinSessionGenerater() StorableSession {
	return NewXinSession()
}

//NewXinRedisHandler return a redis handler use xinSession as generater
func NewXinRedisHandler(store radix.Client, ttl int) *Redis {
	return NewRedisHandler(store, ttl, xinSessionGenerater)
}

//Load implement Handler
func (s *Redis) Load(sessionID string) (Session, bool, error) {
	var raw []byte
	err := s.store.Do(radix.Cmd(&raw, "GET", sessionID))
	if err != nil {
		return nil, false, err
	}
	if len(raw) < 1 {
		return nil, false, nil
	}
	session := s.sessionGenerater()
	err = session.Unmarshal(raw)
	if err != nil {
		return nil, true, err
	}
	return session, true, nil
}

//Save implement Handler
func (s *Redis) Save(sessionID string, session Session) (ttl int, err error) {
	ttl = session.GetTTL()
	if ttl == 0 {
		ttl = s.TTL
	}
	sc, ok := session.(StorableSession)
	if !ok {
		return 0, xin.NewWrapEf("session is not storable")
	}
	if !sc.HasNewContent() {
		// 没有新的内容只做session 的刷新，不创建新内容
		err = s.store.Do(radix.FlatCmd(nil, "EXPIRE", sessionID, ttl))
		if err != nil {
			return 0, err
		}
	} else {
		sBytes, err := sc.Marshal()
		if err != nil {
			return 0, err
		}
		err = s.store.Do(radix.FlatCmd(nil, "SETEX", sessionID, ttl, sBytes))
		if err != nil {
			return 0, err
		}
	}
	return s.TTL, nil
}
