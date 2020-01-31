package sql

import (
	"github.com/lib/pq"
	"github.com/xieqiaoyu/xin"
	"xorm.io/xorm"
)

//XormConfig config support xorm setting
type XormConfig interface {
	Config
	EnableDbLog() bool
}

func newXormEngineHandler(config XormConfig) GenEngineFunc {
	return func(driverName, dataSourceName string) (engine interface{}, err error) {
		e, err := xorm.NewEngine(driverName, dataSourceName)
		if err != nil {
			return nil, xin.NewWrapEf("Fail to new xorm database engine driver [%s] source [%s], Err:%w", driverName, dataSourceName, err)

		}

		logEnable := config.EnableDbLog()
		if logEnable {
			e.ShowSQL(true)
			//engine.Logger().SetLevel(core.LOG_DEBUG)
		}
		return e, nil
	}
}

func closeXormEngine(engine interface{}) error {
	e, ok := engine.(*xorm.Engine)
	if !ok {
		return xin.NewWrapEf("engine is not a *xorm.Engine")
	}
	return e.Close()
}

//XormService xorm engine service
type XormService struct {
	*Service
	config XormConfig
}

//NewXormService NewXormService
func NewXormService(config XormConfig) *XormService {
	return &XormService{
		config:  config,
		Service: NewService(config, newXormEngineHandler(config), closeXormEngine),
	}
}

//Engine load an xorm engine by id
func (s *XormService) Engine(id string) (engine *xorm.Engine, err error) {
	e, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	engine, ok := e.(*xorm.Engine)
	if !ok {
		return nil, xin.NewWrapEf("db id %s is not a *xorm.Engine", id)
	}
	return engine, nil
}

//Session  load session by id if giving inf is nil ,if isNew is true caller  should close session after everything is done
func (s *XormService) Session(id string) (session *xorm.Session, err error) {
	engine, err := s.Engine(id)
	if err != nil {
		return nil, err
	}
	return engine.NewSession(), nil

}

//XormPqStringArray postgres array type support for xorm
type XormPqStringArray []string

//FromDB xorm custom datatype implement
func (a *XormPqStringArray) FromDB(bts []byte) error {
	pqArray := new(pq.StringArray)
	err := pqArray.Scan(bts)
	if err != nil {
		return err
	}
	*a = []string(*pqArray)
	return nil
}

//ToDB xorm custom datatype implement
func (a XormPqStringArray) ToDB() ([]byte, error) {
	pqArray := pq.StringArray(a)
	v, err := pqArray.Value()
	return []byte(v.(string)), err
}
