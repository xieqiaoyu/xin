package handler

import (
	"bytes"
	"encoding/gob"
	"github.com/mediocregopher/radix/v3"
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
	session = &xsession.XinSession{}
	b := bytes.Buffer{}
	b.Write(raw)
	d := gob.NewDecoder(&b)
	err = d.Decode(session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

//NewRedisHandler NewRedisHandler
func NewRedisHandler(store radix.Client, ttl int, marshal SessionMarshalFunc, unmarshal SessionUnmarshalFunc) *Redis {
	if marshal == nil {
		marshal = DefaultSessionMarshal
	}
	if unmarshal == nil {
		unmarshal = DefaultSessionUnmarshal
	}
	return &Redis{
		TTL:       ttl,
		Marshal:   marshal,
		Unmarshal: unmarshal,
		store:     store,
	}
}

//Redis  redis session 控制类
type Redis struct {
	TTL       int
	Marshal   SessionMarshalFunc
	Unmarshal SessionUnmarshalFunc
	store     radix.Client
}

//Load Load
func (s *Redis) Load(sessionID string) (xsession.Session, bool, error) {
	var raw []byte
	err := s.store.Do(radix.Cmd(&raw, "GET", sessionID))
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
	ttl = session.GetTTL()
	if ttl == 0 {
		ttl = s.TTL
	}
	sc, ok := session.(xsession.SessionContent)
	if ok && !sc.HasNewContent() {
		// 没有新的内容只做session 的刷新，不创建新内容
		err = s.store.Do(radix.FlatCmd(nil, "EXPIRE", sessionID, ttl))
		if err != nil {
			return 0, err
		}
	} else {
		sBytes, err := s.Marshal(session)
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
