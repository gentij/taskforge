package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
}

var authToken string

func init() {
	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Login with an API token",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, path, err := loadConfig()
			if err != nil {
				return err
			}

			token := strings.TrimSpace(authToken)
			if token == "" {
				fmt.Print("API token: ")
				reader := bufio.NewReader(os.Stdin)
				line, err := reader.ReadString('\n')
				if err != nil {
					return err
				}
				token = strings.TrimSpace(line)
			}

			if token == "" {
				return fmt.Errorf("token is required")
			}

			cfg.Token = token
			cfg.ServerURL = serverURL

			if err := saveConfig(path, cfg); err != nil {
				return err
			}

			fmt.Printf("Saved credentials to %s\n", path)
			return nil
		},
	}
	loginCmd.Flags().StringVar(&authToken, "token", "", "API token")

	logoutCmd := &cobra.Command{
		Use:   "logout",
		Short: "Clear saved credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, path, err := loadConfig()
			if err != nil {
				return err
			}

			cfg.Token = ""

			if err := saveConfig(path, cfg); err != nil {
				return err
			}

			fmt.Printf("Cleared credentials in %s\n", path)
			return nil
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show auth status",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, path, err := loadConfig()
			if err != nil {
				return err
			}

			server := cfg.ServerURL
			if server == "" {
				server = serverURL
			}

			fmt.Printf("Config: %s\n", path)
			fmt.Printf("Server: %s\n", server)
			fmt.Printf("Token:  %s\n", tokenStatus(cfg.Token))
			return nil
		},
	}

	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(statusCmd)
}

func tokenStatus(token string) string {
	if strings.TrimSpace(token) == "" {
		return "not set"
	}

	return "set"
}
