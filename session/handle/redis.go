package handle

import (
	"bytes"
	"encoding/gob"
	"github.com/mediocregopher/radix/v3"
	"github.com/xieqiaoyu/xin/db/redis"
	xsession "github.com/xieqiaoyu/xin/session"
)

type SessionMarshalFunc func(session xsession.Session) (result []byte, err error)
type SessionUnmarshalFunc func(raw []byte) (session xsession.Session, err error)

func DefaultSessionMarshal(session xsession.Session) ([]byte, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(session)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func DefaultSessionUnmarshal(raw []byte) (session xsession.Session, err error) {
	session = &xsession.MeepoSession{}
	b := bytes.Buffer{}
	b.Write(raw)
	d := gob.NewDecoder(&b)
	err = d.Decode(session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

//NewRedisHandle NewRedisHandle
func NewRedisHandle(ttl int, connectID string, marshal SessionMarshalFunc, unmarshal SessionUnmarshalFunc) *Redis {
	if marshal == nil {
		marshal = DefaultSessionMarshal
	}
	if unmarshal == nil {
		unmarshal = DefaultSessionUnmarshal
	}
	return &Redis{
		TTL:       ttl,
		ConnectID: connectID,
		Marshal:   marshal,
		Unmarshal: unmarshal,
	}
}

//Redis  redis session 控制类
type Redis struct {
	TTL       int
	ConnectID string
	Marshal   SessionMarshalFunc
	Unmarshal SessionUnmarshalFunc
}

//Load Load
func (s *Redis) Load(sessionID string) (xsession.Session, bool, error) {
	store, err := redis.Engine(s.ConnectID)
	if err != nil {
		return nil, false, err
	}
	var raw []byte
	err = store.Do(radix.Cmd(&raw, "GET", sessionID))
	if err != nil {
		return nil, false, err
	}
	if len(raw) < 1 {
		return nil, false, nil
	}
	session, err := s.Unmarshal(raw)
	if err != nil {
		return nil, true, err
	}
	return session, true, nil
}

//Save Save
func (s *Redis) Save(sessionID string, session xsession.Session) (ttl int, err error) {
	store, err := redis.Engine(s.ConnectID)
	if err != nil {
		return 0, err
	}
	ttl = session.GetTTL()
	if ttl == 0 {
		ttl = s.TTL
	}
	sc, ok := session.(xsession.SessionContent)
	if ok && !sc.HasNewContent() {
		// 没有新的内容只做session 的刷新，不创建新内容
		err = store.Do(radix.FlatCmd(nil, "EXPIRE", sessionID, ttl))
		if err != nil {
			return 0, err
		}
	} else {
		sBytes, err := s.Marshal(session)
		if err != nil {
			return 0, err
		}
		err = store.Do(radix.FlatCmd(nil, "SETEX", sessionID, ttl, sBytes))
		if err != nil {
			return 0, err
		}
	}
	return s.TTL, nil
}
