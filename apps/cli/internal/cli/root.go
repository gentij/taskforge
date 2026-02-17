package cli

import (
	"os"

	"github.com/gentij/taskforge/apps/cli/internal/api"
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
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfg, _, err := loadConfig()
		if err != nil {
			return err
		}

		ctx := &Context{
			Config:     cfg,
			Client:     api.NewClient(cfg.ServerURL, cfg.Token),
			OutputJSON: outputJSON,
			Quiet:      quiet,
		}

		cmd.SetContext(WithContext(cmd.Context(), ctx))
		return nil
	},
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
