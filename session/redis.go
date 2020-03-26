package session

import (
	"github.com/go-redis/redis/v7"
	"github.com/xieqiaoyu/xin"
	"time"
)

//StorableSessionGenerater StorableSessionGenerater
type StorableSessionGenerater func() StorableSession

//Redis  redis session handler
type Redis struct {
	TTL              int
	store            *redis.Client
	sessionGenerater StorableSessionGenerater
}

//NewRedisHandler NewRedisHandler
func NewRedisHandler(store *redis.Client, ttl int, generater StorableSessionGenerater) *Redis {
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
func NewXinRedisHandler(store *redis.Client, ttl int) *Redis {
	return NewRedisHandler(store, ttl, xinSessionGenerater)
}

//Load implement Handler
func (s *Redis) Load(sessionID string) (Session, bool, error) {
	raw, err := s.store.Get(sessionID).Bytes()
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
	ttlDuration := time.Duration(ttl) * time.Second

	if !sc.HasNewContent() {
		// 没有新的内容只做session 的刷新，不创建新内容
		_, err = s.store.Expire(sessionID, ttlDuration).Result()
		if err != nil {
			return 0, err
		}
	} else {
		sBytes, err := sc.Marshal()
		if err != nil {
			return 0, err
		}
		_, err = s.store.Set(sessionID, sBytes, ttlDuration).Result()
		if err != nil {
			return 0, err
		}
	}
	return s.TTL, nil
}
