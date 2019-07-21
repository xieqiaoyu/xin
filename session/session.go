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

type MeepoSession struct {
	newContent bool
	Values     map[string]interface{}
	ttl        int
}

func (s *MeepoSession) HasNewContent() bool {
	return s.newContent
}

func (s *MeepoSession) Get(key string) (interface{}, bool) {
	v, exists := s.Values[key]
	return v, exists
}

func (s *MeepoSession) Set(key string, value interface{}) error {
	s.Values[key] = value
	s.newContent = true
	return nil
}

func (s *MeepoSession) Delete(key string) error {
	delete(s.Values, key)
	s.newContent = true
	return nil
}

func (s *MeepoSession) SetTTL(ttl int) error {
	s.ttl = ttl
	return nil
}

func (s *MeepoSession) GetTTL() int {
	return s.ttl
}

func NewSession() Session {
	return &MeepoSession{
		newContent: true,
		Values:     map[string]interface{}{},
	}
}
