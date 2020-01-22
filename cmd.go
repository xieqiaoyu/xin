package xin

import (
	"fmt"

	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

//VerifiableConfig a config interface for config test
type VerifiableConfig interface {
	Verify() error
}

//NewConfigTestCmd return a cobra command for configtest
// command: configtest
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

//NewVersionCmd return a cobra command to print giving version string
// command: version
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

//WaitForQuitSignal block until get  a quit signal (SIGINT,SIGTERM) ,use for graceful stop
func WaitForQuitSignal() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
