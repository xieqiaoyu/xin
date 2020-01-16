package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xieqiaoyu/xin"
)

func NewConfigTestCmd(config *xin.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "configtest",
		Short: "config check",
		Long:  `check config file is ok`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := config.Init(); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("config check pass!")
			}
		},
	}
}
