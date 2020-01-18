package redis

import (
	"github.com/xieqiaoyu/xin"
	"sync"

	"github.com/mediocregopher/radix/v3"
)

type RedisConfig interface {
	GetRedisURI(id string) (string, error)
}

type Service struct {
	instances *sync.Map
	config    RedisConfig
}

func NewService(config RedisConfig) *Service {
	return &Service{
		instances: new(sync.Map),
		config:    config,
	}
}

//Engine 获取redis 连接对象
func (s *Service) Engine(id string) (radix.Client, error) {
	instance, exists := s.instances.Load(id)
	if exists {
		return instance.(radix.Client), nil
	}
	redisURI, err := s.config.GetRedisURI(id)
	if err != nil {
		return nil, xin.NewWrapE(err)
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
