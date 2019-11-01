package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

//ConfigTestCmd ConfigTestCmd
func ConfigTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "configtest",
		Short: "config check",
		Long:  `check config file is ok`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := InitConfig(); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("config check pass!")
			}
		},
	}
}
