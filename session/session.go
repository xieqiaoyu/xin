package session

type Session interface {
	Get(key string) (value interface{}, exists bool)
	Set(key string, value interface{}) error
	Delete(key string) error
	SetTTL(ttl int) error
	GetTTL() int
}

type SessionContent interface {
	HasNewContent() bool
}

type XinSession struct {
	newContent bool
	Values     map[string]interface{}
	ttl        int
}

func (s *XinSession) HasNewContent() bool {
	return s.newContent
}

func (s *XinSession) Get(key string) (interface{}, bool) {
	v, exists := s.Values[key]
	return v, exists
}

func (s *XinSession) Set(key string, value interface{}) error {
	s.Values[key] = value
	s.newContent = true
	return nil
}

func (s *XinSession) Delete(key string) error {
	delete(s.Values, key)
	s.newContent = true
	return nil
}

func (s *XinSession) SetTTL(ttl int) error {
	s.ttl = ttl
	return nil
}

func (s *XinSession) GetTTL() int {
	return s.ttl
}

func NewSession() Session {
	return &XinSession{
		newContent: true,
		Values:     map[string]interface{}{},
	}
}
