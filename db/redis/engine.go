package redis

import (
	"github.com/go-redis/redis/v7"
	"github.com/xieqiaoyu/xin"
	"sync"
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

//NewService create a new redis connect service
func NewService(config Config) *Service {
	return &Service{
		instances: new(sync.Map),
		config:    config,
	}
}

//Engine get redis client by id
func (s *Service) Engine(id string) (*redis.Client, error) {
	instance, exists := s.instances.Load(id)
	if exists {
		return instance.(*redis.Client), nil
	}
	redisURI, err := s.config.GetRedisURI(id)
	if err != nil {
		return nil, xin.NewTracedE(err)
	}
	options, err := redis.ParseURL(redisURI)
	if err != nil {
		return nil, xin.WrapEf(&xin.InternalError{}, "Fail to parse redis [%s] URL:%s", id, err)
	}
	client := redis.NewClient(options)

	instance, loaded := s.instances.LoadOrStore(id, client)
	if loaded {
		client.Close()
	}
	return instance.(*redis.Client), nil
}
