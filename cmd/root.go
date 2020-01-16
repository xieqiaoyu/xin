package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/xieqiaoyu/xin"
	xlog "github.com/xieqiaoyu/xin/log"
)

type SubCmds []*cobra.Command

type RootCmd struct {
	Command *cobra.Command
}

func (c *RootCmd) Execute() {
	if err := c.Command.Execute(); err != nil {
		xlog.WriteError(err.Error())
		os.Exit(1)
	}
}

func NewRootCmd(subcmds SubCmds, fileConfigLoader *xin.FileConfigLoader) *RootCmd {
	cobraCmd := &cobra.Command{
		Use: "anonymous",
	}
	cobraCmd.PersistentFlags().StringVarP(&fileConfigLoader.FileName, "config", "c", "", "specific config file")
	rootcmd := &RootCmd{
		Command: cobraCmd,
	}
	for _, cmd := range subcmds {
		cobraCmd.AddCommand(cmd)
	}
	return rootcmd
}
