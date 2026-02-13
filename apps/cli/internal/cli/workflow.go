package cli

import "github.com/spf13/cobra"

var workflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Manage workflows",
}

func init() {
	workflowCmd.AddCommand(newNotImplementedCmd("list", "List workflows"))
	workflowCmd.AddCommand(newNotImplementedCmd("get", "Get a workflow"))
	workflowCmd.AddCommand(newNotImplementedCmd("create", "Create a workflow"))
	workflowCmd.AddCommand(newNotImplementedCmd("update", "Update a workflow"))
	workflowCmd.AddCommand(newNotImplementedCmd("delete", "Delete a workflow"))
	workflowCmd.AddCommand(newNotImplementedCmd("run", "Run a workflow"))
}
