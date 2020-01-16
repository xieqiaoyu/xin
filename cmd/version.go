package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

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
