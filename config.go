package xin

import (
	"github.com/spf13/viper"
	"sync"
)

var v *viper.Viper
var configOnce sync.Once

//V 获取config viper 对象
func Config() *viper.Viper {
	configOnce.Do(func() {
		v = viper.New()
	})
	return v
}
