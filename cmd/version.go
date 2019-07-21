package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var version = "Unknown"

// SetVersion 设置app 版本号 用于 version 命令显示
func SetVersion(versionString string) {
	version = versionString
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "print the version number",
		Long:  `we also have a version number`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}
}
