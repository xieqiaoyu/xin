package session

import (
	"bytes"
	"encoding/gob"
)

//Session session interface
type Session interface {
	// Get session content by key
	Get(key string) (value interface{}, exists bool)
	// Set session content by key
	Set(key string, value interface{}) error
	// Delete session content by key
	Delete(key string) error
	// SetTTL set session ttl second
	SetTTL(ttl int) error
	// GetTTL get session ttl second
	GetTTL() int
}

//StorableSession storable session interface
type StorableSession interface {
	Session
	// whether Session struct has new content after last generate
	HasNewContent() bool
	// marshal session into byte
	Marshal() (result []byte, err error)
	// unmarshal session from byte
	Unmarshal(raw []byte) (err error)
}

//XinSession a simple implementation of Session interface
type XinSession struct {
	newContent bool
	Values     map[string]interface{}
	ttl        int
}

//HasNewContent implement StorableSession
func (s *XinSession) HasNewContent() bool {
	return s.newContent
}

//Marshal implement StorableSession
func (s *XinSession) Marshal() (result []byte, err error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err = e.Encode(s)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

//Unmarshal implement StorableSession
func (s *XinSession) Unmarshal(raw []byte) (err error) {
	b := bytes.Buffer{}
	b.Write(raw)
	d := gob.NewDecoder(&b)
	err = d.Decode(s)
	if err != nil {
		return err
	}
	s.newContent = false
	return nil
}

//Get implement Session
func (s *XinSession) Get(key string) (interface{}, bool) {
	v, exists := s.Values[key]
	return v, exists
}

//Set implement Session
func (s *XinSession) Set(key string, value interface{}) error {
	s.Values[key] = value
	s.newContent = true
	return nil
}

//Delete implement Session
func (s *XinSession) Delete(key string) error {
	delete(s.Values, key)
	s.newContent = true
	return nil
}

//SetTTL implement Session
func (s *XinSession) SetTTL(ttl int) error {
	s.ttl = ttl
	return nil
}

//GetTTL implement Session
func (s *XinSession) GetTTL() int {
	return s.ttl
}

//NewXinSession create a new xin session
func NewXinSession() *XinSession {
	return &XinSession{
		newContent: true,
		Values:     map[string]interface{}{},
	}
}
