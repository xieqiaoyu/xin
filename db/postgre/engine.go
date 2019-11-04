package postgre

import (
	"fmt"
	"github.com/xieqiaoyu/xin"
	"sync"

	"xorm.io/xorm"
	// 需要执行postgres 包的 init
	_ "github.com/lib/pq"
)

const logEnableKey = "database_enable_log"
const configSourceKey = "database_connections"

var dbInstances *sync.Map

func init() {
	dbInstances = new(sync.Map)
}

// Engine  获取数据库连接对象
func Engine(ids ...string) (*xorm.Engine, error) {
	id := "default"
	if len(ids) > 0 {
		id = ids[0]
	}
	dbInstance, exists := dbInstances.Load(id)
	if exists {
		// 连接已经存在的情况，做一个断言直接返回即可
		return dbInstance.(*xorm.Engine), nil
	}

	conf := xin.Config()
	connectionSourceKey := fmt.Sprintf("%s.%s", configSourceKey, id)
	dbSource := conf.GetString(connectionSourceKey)

	if dbSource == "" {
		return nil, xin.WrapE(&xin.InternalError{}, "Fail to get database source string, please check config key %s in %s", connectionSourceKey, conf.ConfigFileUsed())
	}
	dbInstanceTemp, err := xorm.NewEngine("postgres", dbSource)
	if err != nil {
		return nil, xin.WrapE(&xin.InternalError{}, "Fail to connect database use source string %s, Err:%w", dbSource, err)

	}

	logEnable := conf.GetBool(logEnableKey)
	if logEnable {
		dbInstanceTemp.ShowSQL(true)
		//dbInstanceTemp.Logger().SetLevel(core.LOG_DEBUG)
	}

	dbInstance, loaded := dbInstances.LoadOrStore(id, dbInstanceTemp)
	if loaded {
		// 已经有其他线程打开了数据库连接，本次操作的连接可以关闭
		dbInstanceTemp.Close()
	}

	return dbInstance.(*xorm.Engine), nil
}

//GetOrLoad load session by id if giving inf is nil ,if isNew is true caller  should close session after everything is done
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
	return nil, false, xin.WrapE(&xin.InternalError{}, "Unknown xorm interface type %T", dbInf)

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
