package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd().AddCommand(ConfigTestCmd())
}

//ConfigTestCmd ConfigTestCmd
func ConfigTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "configtest",
		Short: "config check",
		Long:  `check config file is ok`,
		Run: func(cmd *cobra.Command, args []string) {
			InitConfig()
			fmt.Println("config check pass!")
		},
	}
}
