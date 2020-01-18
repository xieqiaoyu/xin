package xin

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	xjsonschema "github.com/xieqiaoyu/xin/util/jsonschema"
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
	err := vc.ReadInConfig()
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

func NewJSONSchemaConfigVerifier(schema string) *JSONSchemaConfigVerifier {
	return &JSONSchemaConfigVerifier{
		Schema: schema,
	}
}

func NewFileConfigLoader(filename, configType string) *FileConfigLoader {
	return &FileConfigLoader{
		FileName:   filename,
		ConfigType: configType,
	}
}

type Config struct {
	loader   ConfigLoader
	verifier ConfigVerifier
	viper    *viper.Viper
}

func NewConfig(configloader ConfigLoader, configVerifier ConfigVerifier) *Config {
	return &Config{
		loader:   configloader,
		verifier: configVerifier,
	}
}

func (c *Config) Init() error {
	v := viper.New()
	if c.loader == nil {
		return fmt.Errorf("Can not init config with a nil config loader")
	}
	err := c.loader.LoadConfig(v)
	if err != nil {
		return err
	}
	// 验证配置文件的内容是否正确
	if c.verifier != nil {
		err := c.verifier.VerfiyConfig(v)
		if err != nil {
			return err
		}
	}
	c.viper = v
	return nil
}

func (c *Config) Verify() error {
	return c.Init()
}

func (c *Config) HttpListen() string {
	return c.viper.GetString("http.listen")
}

func (c *Config) Env() string {
	return c.viper.GetString("env")
}

func (c *Config) EnableDbLog() bool {
	return c.viper.GetBool("database_enable_log")
}

func (c *Config) GetPostgreSource(id string) (string, error) {
	connectionSourceKey := fmt.Sprintf("%s.%s", "database_connections", id)
	dbSource := c.viper.GetString(connectionSourceKey)
	if dbSource == "" {
		return "", WrapEf(&InternalError{}, "Fail to get database source string, please check config key %s in %s", connectionSourceKey, c.viper.ConfigFileUsed())

	}
	return dbSource, nil
}

func (c *Config) GetRedisURI(id string) (string, error) {
	connectionSourceKey := fmt.Sprintf("%s.%s", "redis_connections", id)
	redisURI := c.viper.GetString(connectionSourceKey)
	if redisURI == "" {
		return "", WrapEf(&InternalError{}, "Fail to get redis URI,pleas check config key %s in %s", connectionSourceKey, c.viper.ConfigFileUsed())
	}
	return redisURI, nil
}

func (c *Config) Viper() *viper.Viper {
	return c.viper
}
