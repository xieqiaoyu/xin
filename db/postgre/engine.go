package postgre

import (
	"github.com/xieqiaoyu/xin"
	"sync"

	"xorm.io/xorm"
	// 需要执行postgres 包的 init
	_ "github.com/lib/pq"
)

type PostgreConfig interface {
	EnableDbLog() bool
	GetPostgreSource(id string) (string, error)
}

type Service struct {
	instances *sync.Map
	config    PostgreConfig
}

func NewService(config *xin.Config) *Service {
	return &Service{
		instances: new(sync.Map),
		config:    config,
	}
}

func (s *Service) Engine(id string) (*xorm.Engine, error) {
	dbInstance, exists := s.instances.Load(id)
	if exists {
		// 连接已经存在的情况，做一个断言直接返回即可
		return dbInstance.(*xorm.Engine), nil
	}

	dbSource, err := s.config.GetPostgreSource(id)
	if err != nil {
		return nil, xin.NewWrapE(err)
	}

	dbInstanceTemp, err := xorm.NewEngine("postgres", dbSource)
	if err != nil {
		return nil, xin.WrapEf(&xin.InternalError{}, "Fail to connect database use source string %s, Err:%w", dbSource, err)

	}

	logEnable := s.config.EnableDbLog()
	if logEnable {
		dbInstanceTemp.ShowSQL(true)
		//dbInstanceTemp.Logger().SetLevel(core.LOG_DEBUG)
	}

	dbInstance, loaded := s.instances.LoadOrStore(id, dbInstanceTemp)
	if loaded {
		// 已经有其他线程打开了数据库连接，本次操作的连接可以关闭
		dbInstanceTemp.Close()
	}

	return dbInstance.(*xorm.Engine), nil
}

func (s *Service) Session(id string) (session *xorm.Session, err error) {
	engine, err := s.Engine(id)
	if err != nil {
		return nil, err
	}
	return engine.NewSession(), nil

}
