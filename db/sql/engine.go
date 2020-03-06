package sql

import (
	"github.com/xieqiaoyu/xin"
	"sync"
)

//Config config provide sql connection setting
type Config interface {
	GetSQLSource(id string) (driver string, source string, err error)
}

//GenEngineFunc  function to generate an engine
type GenEngineFunc func(driverName, dataSourceName string) (engine interface{}, err error)

//CloseEngineFunc function to close an engine
type CloseEngineFunc func(engine interface{}) error

//Service common sql connect service
type Service struct {
	instances   *sync.Map
	config      Config
	genHandle   GenEngineFunc
	closeHandle CloseEngineFunc
}

//NewService create a new sql connect service
func NewService(config Config, genHandle GenEngineFunc, closeHandle CloseEngineFunc) *Service {
	return &Service{
		instances:   new(sync.Map),
		config:      config,
		genHandle:   genHandle,
		closeHandle: closeHandle,
	}
}

// Get  get sql connect engine
func (s *Service) Get(id string) (interface{}, error) {
	if s.genHandle == nil {
		return nil, xin.NewTracedEf("genHandle can not be nil")
	}
	dbInstance, exists := s.instances.Load(id)
	if exists {
		return dbInstance, nil
	}

	sqlDriver, sqlSource, err := s.config.GetSQLSource(id)
	if err != nil {
		return nil, xin.NewTracedE(err)
	}

	dbInstanceTemp, err := s.genHandle(sqlDriver, sqlSource)
	if err != nil {
		return nil, xin.WrapEf(&xin.InternalError{}, "create engine Err:%w", err)

	}

	dbInstance, loaded := s.instances.LoadOrStore(id, dbInstanceTemp)
	if loaded {
		// another routine has already opened the connection, just close ours
		if s.closeHandle != nil {
			s.closeHandle(dbInstanceTemp)
		}
	}

	return dbInstance, nil
}
