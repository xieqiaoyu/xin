package redis

import (
	"github.com/xieqiaoyu/xin"
	"sync"

	"github.com/mediocregopher/radix/v3"
)

//Config config provide redis connection setting
type Config interface {
	GetRedisURI(id string) (string, error)
}

//Service redis connect service
type Service struct {
	instances *sync.Map
	config    Config
}

//NewService create a new radis connect service
func NewService(config Config) *Service {
	return &Service{
		instances: new(sync.Map),
		config:    config,
	}
}

//Engine get radix client by id
func (s *Service) Engine(id string) (radix.Client, error) {
	instance, exists := s.instances.Load(id)
	if exists {
		return instance.(radix.Client), nil
	}
	redisURI, err := s.config.GetRedisURI(id)
	if err != nil {
		return nil, xin.NewTracedE(err)
	}

	redisPool, err := radix.NewPool("tcp", redisURI, 10)
	if err != nil {
		return nil, xin.WrapEf(&xin.InternalError{}, "Fail to create redis connect pool %w", err)
	}
	instance, loaded := s.instances.LoadOrStore(id, redisPool)
	if loaded {
		redisPool.Close()
	}
	return instance.(radix.Client), nil
}
