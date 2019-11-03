package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xieqiaoyu/xin"
)

//ConfigTestCmd ConfigTestCmd
func ConfigTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "configtest",
		Short: "config check",
		Long:  `check config file is ok`,
		Run: func(cmd *cobra.Command, args []string) {
			if ConfigFileToUse != "" {
				xin.SetConfigFile(ConfigFileToUse)
			}
			if err := xin.LoadConfig(); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("config check pass!")
			}
		},
	}
}
