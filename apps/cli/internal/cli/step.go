package cli

import "github.com/spf13/cobra"

var stepCmd = &cobra.Command{
	Use:   "step",
	Short: "Manage step runs",
}

func init() {
	stepCmd.AddCommand(newNotImplementedCmd("list", "List step runs"))
	stepCmd.AddCommand(newNotImplementedCmd("get", "Get a step run"))
}
