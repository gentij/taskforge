package cli

import "github.com/spf13/cobra"

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Manage workflow runs",
}

func init() {
	runCmd.AddCommand(newNotImplementedCmd("list", "List workflow runs"))
	runCmd.AddCommand(newNotImplementedCmd("get", "Get a workflow run"))
	runCmd.AddCommand(newNotImplementedCmd("cancel", "Cancel a workflow run"))
}
