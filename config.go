package xin

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	xjsonschema "github.com/xieqiaoyu/xin/util/jsonschema"
)

//ConfigLoader load config interface
type ConfigLoader interface {
	LoadConfig(vc *viper.Viper) error
}

//ConfigVerifier verify config interface
type ConfigVerifier interface {
	VerfiyConfig(vc *viper.Viper) error
}

//FileConfigLoader load config from file system
type FileConfigLoader struct {
	FileName   string
	ConfigType string
}

//LoadConfig ConfigLoader interface
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

//JSONSchemaConfigVerifier verfiy config by jsonschema
type JSONSchemaConfigVerifier struct {
	Schema string
}

//VerfiyConfig ConfigVerifier interface
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

//NewJSONSchemaConfigVerifier create a new jsonschema config verifier
func NewJSONSchemaConfigVerifier(schema string) *JSONSchemaConfigVerifier {
	return &JSONSchemaConfigVerifier{
		Schema: schema,
	}
}

//NewFileConfigLoader create a new file config loader
func NewFileConfigLoader(filename, configType string) *FileConfigLoader {
	return &FileConfigLoader{
		FileName:   filename,
		ConfigType: configType,
	}
}

//Config Config
type Config struct {
	loader   ConfigLoader
	verifier ConfigVerifier
	viper    *viper.Viper
}

//NewConfig create a new config
func NewConfig(configloader ConfigLoader, configVerifier ConfigVerifier) *Config {
	return &Config{
		loader:   configloader,
		verifier: configVerifier,
	}
}

//Init init config,load config and verfiy ,this method must be called before other method
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

//Verify verify config
func (c *Config) Verify() error {
	return c.Init()
}

//HTTPListen get http listen string
func (c *Config) HTTPListen() string {
	return c.viper.GetString("http.listen")
}

//Env get env string
func (c *Config) Env() string {
	return c.viper.GetString("env")
}

//EnableDbLog get config for enable db Log
func (c *Config) EnableDbLog() bool {
	return c.viper.GetBool("database_enable_log")
}

//GetPostgreSource get source string for postgresql
func (c *Config) GetPostgreSource(id string) (string, error) {
	connectionSourceKey := fmt.Sprintf("%s.%s", "database_connections", id)
	dbSource := c.viper.GetString(connectionSourceKey)
	if dbSource == "" {
		return "", WrapEf(&InternalError{}, "Fail to get database source string, please check config key %s in %s", connectionSourceKey, c.viper.ConfigFileUsed())

	}
	return dbSource, nil
}

//GetRedisURI get redis connect string
func (c *Config) GetRedisURI(id string) (string, error) {
	connectionSourceKey := fmt.Sprintf("%s.%s", "redis_connections", id)
	redisURI := c.viper.GetString(connectionSourceKey)
	if redisURI == "" {
		return "", WrapEf(&InternalError{}, "Fail to get redis URI,pleas check config key %s in %s", connectionSourceKey, c.viper.ConfigFileUsed())
	}
	return redisURI, nil
}

//GrpcListen get grpc listen info
func (c *Config) GrpcListen() (network, address string) {
	return c.viper.GetString("grpc.network"), c.viper.GetString("grpc.listen")
}

//Viper Get viper instance of config
func (c *Config) Viper() *viper.Viper {
	return c.viper
}
