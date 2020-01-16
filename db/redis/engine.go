package redis

import (
	"fmt"
	"github.com/xieqiaoyu/xin"
	xlog "github.com/xieqiaoyu/xin/log"
	"sync"

	"github.com/mediocregopher/radix/v3"
)

const configSourceKey = "redis_connections"

type Service struct {
	instances *sync.Map
	config    *xin.Config
}

func NewService(config *xin.Config) *Service {
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
	v := s.config.Viper()
	connectionSourceKey := fmt.Sprintf("%s.%s", configSourceKey, id)
	redisURI := v.GetString(connectionSourceKey)
	if redisURI == "" {
		xlog.WithTag("CONFIG").WriteCritical("Fail to get redis uri string, please check config key %s in %s", configSourceKey, v.ConfigFileUsed())
		return nil, xin.WrapEf(&xin.InternalError{}, "Fail to get redis URI")
	}

	redisPool, err := radix.NewPool("tcp", redisURI, 10)
	if err != nil {
		xlog.WithTag("REDIS").WriteCritical("Fail to connect redis use uri %s, Err:%s", redisURI, err)
		return nil, xin.WrapEf(&xin.InternalError{}, "Fail to create redis connect pool")
	}
	instance, loaded := s.instances.LoadOrStore(id, redisPool)
	if loaded {
		redisPool.Close()
	}
	return instance.(radix.Client), nil
}
