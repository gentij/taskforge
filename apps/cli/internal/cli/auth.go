package cli

import "github.com/spf13/cobra"

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
}

func init() {
	authCmd.AddCommand(newNotImplementedCmd("login", "Login with an API token"))
	authCmd.AddCommand(newNotImplementedCmd("logout", "Clear saved credentials"))
	authCmd.AddCommand(newNotImplementedCmd("status", "Show auth status"))
}
