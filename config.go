package xin

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	xjsonschema "github.com/xieqiaoyu/xin/util/jsonschema"
	"sync"
)

type ConfigLoader interface {
	LoadConfig(vc *viper.Viper) error
}

type ConfigVerifier interface {
	VerfiyConfig(vc *viper.Viper) error
}

type FileConfigLoader struct {
	FileName   string
	ConfigType string
}

func (l *FileConfigLoader) LoadConfig(vc *viper.Viper) error {
	if l.ConfigType != "" {
		vc.SetConfigType(l.ConfigType)
	} else {
		vc.SetConfigType("toml")
	}
	if l.FileName != "" {
		vc.SetConfigFile(l.FileName)
	} else {
		vc.AddConfigPath(".")
		vc.SetConfigName("config")
	}
	err := v.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		return fmt.Errorf("Fatal error load config file: %w", err)
	}
	return nil
}

type JSONSchemaConfigVerifier struct {
	Schema string
}

func (jv JSONSchemaConfigVerifier) VerfiyConfig(vc *viper.Viper) error {
	config := make(map[string]interface{})
	err := vc.Unmarshal(&config)
	if err != nil {
		return fmt.Errorf("Unmarshal config fail:: %w", err)
	}
	configString, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("Marshal config fail:%w", err)
	}
	if pass, err := xjsonschema.ValidJSONString(string(configString), jv.Schema); !pass {
		return fmt.Errorf("Unsatisfy config :%w", err)
	}
	return nil
}

var (
	v              *viper.Viper
	configOnce     sync.Once
	configLoader   ConfigLoader
	configVerifier ConfigVerifier
)

//V 获取config viper 对象
func Config() *viper.Viper {
	configOnce.Do(func() {
		v = viper.New()
	})
	return v
}

func SetConfigLoader(l ConfigLoader) {
	configLoader = l
}

func SetConfigFile(filename, configType string) {
	SetConfigLoader(&FileConfigLoader{
		FileName:   filename,
		ConfigType: configType,
	})
}

//VerifyConfigBySchema 开启配置文件校验使用指定JSONschema 来校验config配置是否符合要求
func VerifyConfigBySchema(schema string) {
	configVerifier = &JSONSchemaConfigVerifier{
		Schema: schema,
	}
}

//LoadConfig load Config ,this func should be called before use config
func LoadConfig() error {
	configOnce.Do(func() {
		v = viper.New()
	})
	if configLoader == nil {
		configLoader = &FileConfigLoader{}
	}
	err := configLoader.LoadConfig(v)
	if err != nil {
		return err
	}
	// 验证配置文件的内容是否正确
	if configVerifier != nil {
		err := configVerifier.VerfiyConfig(v)
		if err != nil {
			return err
		}
	}
	return nil
}
