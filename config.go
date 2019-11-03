package xin

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	xjsonschema "github.com/xieqiaoyu/xin/util/jsonschema"
	"sync"
)

type LoadConfigHandle func(vc *viper.Viper) error

func init() {
	loadConfigAct = defaultLoadConfigAct
	ConfigType = "toml"
}

var (
	v             *viper.Viper
	configOnce    sync.Once
	configFile    string
	configSchema  string
	verifyConfig  bool // default no verify config
	loadConfigAct LoadConfigHandle
	ConfigType    string
)

//V 获取config viper 对象
func Config() *viper.Viper {
	configOnce.Do(func() {
		v = viper.New()
	})
	return v
}

func SetConfigFile(fileName string) {
	configFile = fileName
}

func SetConfigLoadAct(act LoadConfigHandle) {
	loadConfigAct = act
}

//VerifyConfigBySchema 开启配置文件校验使用指定JSONschema 来校验config配置是否符合要求
func VerifyConfigBySchema(schema string) {
	configSchema = schema
	verifyConfig = true
}

//LoadConfig load Config ,this func should be called before use config
func LoadConfig() error {
	configOnce.Do(func() {
		v = viper.New()
	})
	err := loadConfigAct(v)
	if err != nil {
		return err
	}
	// 验证配置文件的内容是否正确
	if verifyConfig {
		config := make(map[string]interface{})
		err = v.Unmarshal(&config)
		if err != nil {
			return fmt.Errorf("Unmarshal config fail:: %w", err)
		}
		configString, err := json.Marshal(config)
		if err != nil {
			return fmt.Errorf("Marshal config fail:%w", err)
		}
		if pass, err := xjsonschema.ValidJSONString(string(configString), configSchema); !pass {
			return fmt.Errorf("Unsatisfy config :%w", err)
		}
	}
	return nil
}

func defaultLoadConfigAct(vc *viper.Viper) error {
	if configFile != "" {
		vc.SetConfigFile(configFile)
	} else {
		vc.AddConfigPath(".")
		vc.SetConfigName("config")
	}
	vc.SetConfigType(ConfigType)
	err := v.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		return fmt.Errorf("Fatal error load config file: %w", err)
	}
	return nil
}
