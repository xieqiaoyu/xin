package sql

import (
	"fmt"
	"github.com/xieqiaoyu/xin"
	"sync"

	"xorm.io/xorm"
)

const logEnableKey = "sql_enable_log"
const configSourceKey = "sql_connections.%s.source"
const configDriverKey = "sql_connections.%s.driver"

var dbInstances *sync.Map

func init() {
	dbInstances = new(sync.Map)
}

// Engine  get connect engine
func Engine(ids ...string) (*xorm.Engine, error) {
	id := "default"
	if len(ids) > 0 {
		id = ids[0]
	}
	dbInstance, exists := dbInstances.Load(id)
	if exists {
		return dbInstance.(*xorm.Engine), nil
	}

	conf := xin.Config()
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
	dbInstanceTemp, err := xorm.NewEngine(sqlDriver, sqlSource)
	if err != nil {
		return nil, xin.WrapEf(&xin.InternalError{}, "Fail to connect database use source string %s, Err:%w", sqlSource, err)

	}

	logEnable := conf.GetBool(logEnableKey)
	if logEnable {
		dbInstanceTemp.ShowSQL(true)
		//dbInstanceTemp.Logger().SetLevel(core.LOG_DEBUG)
	}

	dbInstance, loaded := dbInstances.LoadOrStore(id, dbInstanceTemp)
	if loaded {
		// another routine has already opened the connection, just close ours
		dbInstanceTemp.Close()
	}

	return dbInstance.(*xorm.Engine), nil
}

//GetOrLoad load session by id if giving inf is nil ,if `isNew` is true, caller  should close session after everything is done
func Session(id string, dbInf xorm.Interface) (session *xorm.Session, isNew bool, err error) {
	if dbInf == nil {
		engine, err := Engine(id)
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

//Close Close
func Close() {
	dbInstances.Range(func(id, dbInstance interface{}) bool {
		dbE := dbInstance.(*xorm.Engine)
		err := dbE.Close()
		if err != nil {
			return false
		}
		return true
	})
}
