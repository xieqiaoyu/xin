package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xieqiaoyu/xin"
)

type VerifiableConfig interface {
	Verify() error
}

func NewConfigTestCmd(config VerifiableConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "configtest",
		Short: "config check",
		Long:  `check config file is ok`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := config.Verify(); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("config check pass!")
			}
		},
	}
}
