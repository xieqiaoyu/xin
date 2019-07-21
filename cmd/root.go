package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"encoding/json"
	"github.com/xieqiaoyu/xin"
	xlog "github.com/xieqiaoyu/xin/log"
	xjsonschema "github.com/xieqiaoyu/xin/util/jsonschema"
)

var (
	rootCmd = &cobra.Command{
		Use:   "anonymous",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	cfgFile      string
	configSchema string
	verifyConfig bool // default no verify config
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "specific config file")
}

//RootCmd return root commend
func RootCmd() *cobra.Command {
	return rootCmd
}

//VerifyConfigBySchema 开启配置文件校验使用指定JSONschema 来校验config配置是否符合要求
func VerifyConfigBySchema(schema string) {
	configSchema = schema
	verifyConfig = true
}

//InitConfig 初始化配置文件逻辑,如果读取配置文件失败会报错,这个函数需要手动调用，因为不是所有的命令都需要一个配置文件
func InitConfig() {
	viper := xin.Config()
	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		xlog.WriteError("Fatal error load config file: %s", err)
		os.Exit(2)
	}
	// 验证配置文件的内容是否正确
	if verifyConfig {
		config := make(map[string]interface{})
		err = viper.Unmarshal(&config)
		if err != nil {
			xlog.WriteError("Unmarshal config fail:%s", err)
			os.Exit(2)
		}
		configString, err := json.Marshal(config)
		if err != nil {
			xlog.WriteError("Marshal config fail:%s", err)
			os.Exit(2)
		}
		if pass, err := xjsonschema.ValidJSONString(string(configString), configSchema); !pass {
			xlog.WriteError("Unsatisfy config :%s", err)
			os.Exit(2)
		}
	}
}

// Execute Execute
func Execute() {
	rootCmd.AddCommand(versionCmd())
	if err := rootCmd.Execute(); err != nil {
		xlog.WriteError(err.Error())
		os.Exit(1)
	}
}
