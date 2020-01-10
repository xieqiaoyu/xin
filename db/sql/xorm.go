package sql

import (
	"github.com/lib/pq"
	"github.com/xieqiaoyu/xin"
	"xorm.io/xorm"
)

const logEnableKey = "sql_enable_log"

func newXormEngine(driverName, dataSourceName string) (engine interface{}, err error) {
	e, err := xorm.NewEngine(driverName, dataSourceName)
	if err != nil {
		return nil, xin.NewWrapEf("Fail to new xorm database engine driver [%s] source [%s], Err:%w", driverName, dataSourceName, err)

	}

	logEnable := xin.Config().GetBool(logEnableKey)
	if logEnable {
		e.ShowSQL(true)
		//engine.Logger().SetLevel(core.LOG_DEBUG)
	}
	return e, nil
}

func closeXormEngine(engine interface{}) error {
	e, ok := engine.(*xorm.Engine)
	if !ok {
		return xin.NewWrapEf("engine is not a *xorm.Engine")
	}
	return e.Close()
}

//XormEngine load a
func XormEngine(id string) (engine *xorm.Engine, err error) {
	e, err := Engine(id, newXormEngine, closeXormEngine)
	if err != nil {
		return nil, err
	}
	engine, ok := e.(*xorm.Engine)
	if !ok {
		return nil, xin.NewWrapEf("db id %s is not a *xorm.Engine", id)
	}
	return engine, nil
}

//XormSession  load session by id if giving inf is nil ,if isNew is true caller  should close session after everything is done
func XormSession(id string, dbInf xorm.Interface) (session *xorm.Session, isNew bool, err error) {
	if dbInf == nil {
		engine, err := XormEngine(id)
		if err != nil {
			return nil, false, err
		}
		return engine.NewSession(), true, nil
	}
	switch i := dbInf.(type) {
	case *xorm.Engine:
		return i.NewSession(), true, nil
	case *xorm.Session:
		return i, false, nil
	}
	return nil, false, xin.WrapEf(&xin.InternalError{}, "Unknown xorm interface type %T", dbInf)

}

//XormPqStringArray postgres array type support for xorm
type XormPqStringArray []string

func (a *XormPqStringArray) FromDB(bts []byte) error {
	pqArray := new(pq.StringArray)
	err := pqArray.Scan(bts)
	if err != nil {
		return err
	}
	*a = []string(*pqArray)
	return nil
}

func (a XormPqStringArray) ToDB() ([]byte, error) {
	pqArray := pq.StringArray(a)
	v, err := pqArray.Value()
	return []byte(v.(string)), err
}
