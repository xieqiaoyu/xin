package sql

import (
	"fmt"
	"github.com/xieqiaoyu/xin"
	"sync"
)

const configSourceKey = "sql_connections.%s.source"
const configDriverKey = "sql_connections.%s.driver"

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

//GenEngineFunc  function to generate an engine
type GenEngineFunc func(driverName, dataSourceName string) (engine interface{}, err error)

//CloseEngineFunc function to close an engine
type CloseEngineFunc func(engine interface{}) error

// Engine  get connect engine
func (s *Service) Engine(id string, genHandle GenEngineFunc, closeHandle CloseEngineFunc) (interface{}, error) {
	if genHandle == nil {
		return nil, xin.NewWrapEf("genHandle can not be nil")
	}
	dbInstance, exists := s.instances.Load(id)
	if exists {
		return dbInstance, nil
	}

	conf := s.config.Viper()
	connectionDriverKey := fmt.Sprintf(configDriverKey, id)

	sqlDriver := conf.GetString(connectionDriverKey)
	if sqlDriver == "" {
		return nil, xin.WrapEf(&xin.InternalError{}, "Fail to get sql driver string, please check config key %s in %s", configDriverKey, conf.ConfigFileUsed())
	}
	connectionSourceKey := fmt.Sprintf(configSourceKey, id)
	sqlSource := conf.GetString(connectionSourceKey)

	if sqlSource == "" {
		return nil, xin.WrapEf(&xin.InternalError{}, "Fail to get sql source string, please check config key %s in %s", connectionSourceKey, conf.ConfigFileUsed())
	}
	dbInstanceTemp, err := genHandle(sqlDriver, sqlSource)
	if err != nil {
		return nil, xin.WrapEf(&xin.InternalError{}, "create engine Err:%w", err)

	}

	dbInstance, loaded := s.instances.LoadOrStore(id, dbInstanceTemp)
	if loaded {
		// another routine has already opened the connection, just close ours
		if closeHandle != nil {
			closeHandle(dbInstanceTemp)
		}
	}

	return dbInstance, nil
}
