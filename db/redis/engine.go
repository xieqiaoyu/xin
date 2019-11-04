package redis

import (
	"fmt"
	"github.com/xieqiaoyu/xin"
	xlog "github.com/xieqiaoyu/xin/log"
	"sync"

	"github.com/mediocregopher/radix/v3"
)

const configSourceKey = "redis_connections"

var instances *sync.Map

func init() {
	instances = new(sync.Map)
}

//Engine 获取redis 连接对象
func Engine(ids ...string) (radix.Client, error) {
	id := "default"
	if len(ids) > 0 {
		id = ids[0]
	}
	instance, exists := instances.Load(id)
	if exists {
		return instance.(radix.Client), nil
	}
	v := xin.Config()
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
	instance, loaded := instances.LoadOrStore(id, redisPool)
	if loaded {
		redisPool.Close()
	}
	return instance.(radix.Client), nil
}
