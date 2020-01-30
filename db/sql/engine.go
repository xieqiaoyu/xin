package sql

import (
	"github.com/xieqiaoyu/xin"
	"sync"
)

const configSourceKey = "sql_connections.%s.source"
const configDriverKey = "sql_connections.%s.driver"

//Config config provide sql connection setting
type Config interface {
	EnableDbLog() bool
	GetSQLSource(id string) (driver string, source string, err error)
}

//Service common sql connect service
type Service struct {
	instances *sync.Map
	config    Config
}

//NewService create a new sql connect service
func NewService(config Config) *Service {
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

	sqlDriver, sqlSource, err := s.config.GetSQLSource(id)
	if err != nil {
		return nil, xin.NewWrapE(err)
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
