package cli

import (
	"os"

	"github.com/spf13/cobra"
)

const defaultServerURL = "http://localhost:3000/v1/api"

var (
	configPath string
	serverURL  string
	outputJSON bool
	quiet      bool
)

var rootCmd = &cobra.Command{
	Use:   "taskforge",
	Short: "Taskforge CLI",
	Long:  "Taskforge CLI for managing workflows, runs, triggers, and secrets.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to config file")
	rootCmd.PersistentFlags().StringVar(&serverURL, "server", defaultServerURL, "API server URL")
	rootCmd.PersistentFlags().BoolVar(&outputJSON, "json", false, "Output JSON")
	rootCmd.PersistentFlags().BoolVar(&quiet, "quiet", false, "Minimal output")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(workflowCmd)
	rootCmd.AddCommand(triggerCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(stepCmd)
	rootCmd.AddCommand(secretCmd)
}
