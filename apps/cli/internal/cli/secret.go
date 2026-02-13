package cli

import "github.com/spf13/cobra"

var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Manage secrets",
}

func init() {
	secretCmd.AddCommand(newNotImplementedCmd("list", "List secrets"))
	secretCmd.AddCommand(newNotImplementedCmd("get", "Get a secret"))
	secretCmd.AddCommand(newNotImplementedCmd("create", "Create a secret"))
	secretCmd.AddCommand(newNotImplementedCmd("update", "Update a secret"))
	secretCmd.AddCommand(newNotImplementedCmd("delete", "Delete a secret"))
}
