package xin

import (
	"fmt"

	"github.com/spf13/cobra"
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

func NewVersionCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "print the version number",
		Long:  `we also have a version number`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}
}
