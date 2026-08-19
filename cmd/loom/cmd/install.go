package cmd

import "github.com/spf13/cobra"

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the framework to detected AI platforms",
	RunE:  runStub,
}
