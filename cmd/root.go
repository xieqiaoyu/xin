package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/xieqiaoyu/xin"
	xlog "github.com/xieqiaoyu/xin/log"
)

var (
	rootCmd = &cobra.Command{
		Use: "anonymous",
	}
	ConfigFileToUse string
	ConfigType      = "toml"
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&ConfigFileToUse, "config", "c", "", "specific config file")
}

//RootCmd return root commend
func RootCmd() *cobra.Command {
	return rootCmd
}

// Execute Execute
func Execute() {
	rootCmd.AddCommand(versionCmd())
	rootCmd.AddCommand(ConfigTestCmd())
	if err := rootCmd.Execute(); err != nil {
		xlog.WriteError(err.Error())
		os.Exit(1)
	}
}

func ConfigInit() error {
	if ConfigFileToUse != "" {
		xin.SetConfigFile(ConfigFileToUse, ConfigType)
	}
	return xin.LoadConfig()
}
