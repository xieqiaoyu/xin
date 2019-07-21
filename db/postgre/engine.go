package postgre

import (
	"fmt"
	"github.com/xieqiaoyu/xin"
	"sync"

	"github.com/go-xorm/xorm"
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
		return nil, fmt.Errorf("Fail to get database source string, please check config key %s in %s", connectionSourceKey, conf.ConfigFileUsed())
	}
	dbInstanceTemp, err := xorm.NewEngine("postgres", dbSource)
	if err != nil {
		return nil, fmt.Errorf("Fail to connect database use source string %s, Err:%s", dbSource, err)

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
