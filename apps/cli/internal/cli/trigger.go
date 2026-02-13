package cli

import "github.com/spf13/cobra"

var triggerCmd = &cobra.Command{
	Use:   "trigger",
	Short: "Manage triggers",
}

func init() {
	triggerCmd.AddCommand(newNotImplementedCmd("list", "List triggers"))
	triggerCmd.AddCommand(newNotImplementedCmd("get", "Get a trigger"))
	triggerCmd.AddCommand(newNotImplementedCmd("create", "Create a trigger"))
	triggerCmd.AddCommand(newNotImplementedCmd("update", "Update a trigger"))
	triggerCmd.AddCommand(newNotImplementedCmd("delete", "Delete a trigger"))
}
