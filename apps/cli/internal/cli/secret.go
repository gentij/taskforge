package cli

import (
	"fmt"
	"net/url"
	"os"

	"github.com/gentij/taskforge/apps/cli/internal/api"
	"github.com/gentij/taskforge/apps/cli/internal/output"
	"github.com/spf13/cobra"
)

type secretItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

type secretListResponse struct {
	Items      []secretItem   `json:"items"`
	Pagination paginationMeta `json:"pagination"`
}

var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Manage secrets",
}

var secretListPage int
var secretListPageSize int
var secretCreateName string
var secretCreateValue string
var secretCreateDescription string
var secretUpdateName string
var secretUpdateValue string
var secretUpdateDescription string

func init() {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List secrets",
		RunE:  secretList,
	}
	listCmd.Flags().IntVar(&secretListPage, "page", 1, "Page number")
	listCmd.Flags().IntVar(&secretListPageSize, "page-size", 25, "Page size")

	getCmd := &cobra.Command{
		Use:   "get <secret-id>",
		Short: "Get a secret",
		Args:  cobra.ExactArgs(1),
		RunE:  secretGet,
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a secret",
		RunE:  secretCreate,
	}
	createCmd.Flags().StringVar(&secretCreateName, "name", "", "Secret name")
	createCmd.Flags().StringVar(&secretCreateValue, "value", "", "Secret value")
	createCmd.Flags().StringVar(&secretCreateDescription, "description", "", "Secret description")
	_ = createCmd.MarkFlagRequired("name")
	_ = createCmd.MarkFlagRequired("value")

	updateCmd := &cobra.Command{
		Use:   "update <secret-id>",
		Short: "Update a secret",
		Args:  cobra.ExactArgs(1),
		RunE:  secretUpdate,
	}
	updateCmd.Flags().StringVar(&secretUpdateName, "name", "", "Secret name")
	updateCmd.Flags().StringVar(&secretUpdateValue, "value", "", "Secret value")
	updateCmd.Flags().StringVar(&secretUpdateDescription, "description", "", "Secret description")

	deleteCmd := &cobra.Command{
		Use:   "delete <secret-id>",
		Short: "Delete a secret",
		Args:  cobra.ExactArgs(1),
		RunE:  secretDelete,
	}

	secretCmd.AddCommand(listCmd)
	secretCmd.AddCommand(getCmd)
	secretCmd.AddCommand(createCmd)
	secretCmd.AddCommand(updateCmd)
	secretCmd.AddCommand(deleteCmd)
}

func secretList(cmd *cobra.Command, args []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.ServerURL, cfg.Token)
	query := url.Values{}
	query.Set("page", fmt.Sprintf("%d", secretListPage))
	query.Set("pageSize", fmt.Sprintf("%d", secretListPageSize))

	var result secretListResponse
	if err := client.GetJSON("/secrets?"+query.Encode(), &result); err != nil {
		return err
	}

	if outputJSON {
		return output.PrintJSON(result)
	}

	if quiet {
		for _, item := range result.Items {
			fmt.Fprintln(os.Stdout, item.ID)
		}
		return nil
	}

	w := output.NewTableWriter()
	fmt.Fprintln(w, "ID\tNAME\tCREATED\tUPDATED")
	for _, item := range result.Items {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", item.ID, item.Name, item.CreatedAt, item.UpdatedAt)
	}
	return w.Flush()
}

func secretGet(cmd *cobra.Command, args []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.ServerURL, cfg.Token)
	var result secretItem
	if err := client.GetJSON("/secrets/"+args[0], &result); err != nil {
		return err
	}

	return printSecret(result)
}

func secretCreate(cmd *cobra.Command, args []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	payload := map[string]any{
		"name":        secretCreateName,
		"value":       secretCreateValue,
		"description": secretCreateDescription,
	}

	client := api.NewClient(cfg.ServerURL, cfg.Token)
	var result secretItem
	if err := client.PostJSON("/secrets", payload, &result); err != nil {
		return err
	}

	return printSecret(result)
}

func secretUpdate(cmd *cobra.Command, args []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	patch := map[string]any{}
	if secretUpdateName != "" {
		patch["name"] = secretUpdateName
	}
	if secretUpdateValue != "" {
		patch["value"] = secretUpdateValue
	}
	if cmd.Flags().Changed("description") {
		patch["description"] = secretUpdateDescription
	}
	if len(patch) == 0 {
		return fmt.Errorf("no fields to update")
	}

	client := api.NewClient(cfg.ServerURL, cfg.Token)
	var result secretItem
	if err := client.PatchJSON("/secrets/"+args[0], patch, &result); err != nil {
		return err
	}

	return printSecret(result)
}

func secretDelete(cmd *cobra.Command, args []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.ServerURL, cfg.Token)
	var result secretItem
	if err := client.DeleteJSON("/secrets/"+args[0], &result); err != nil {
		return err
	}

	return printSecret(result)
}

func printSecret(result secretItem) error {
	if outputJSON {
		return output.PrintJSON(result)
	}
	if quiet {
		fmt.Fprintln(os.Stdout, result.ID)
		return nil
	}

	w := output.NewTableWriter()
	fmt.Fprintln(w, "FIELD\tVALUE")
	fmt.Fprintf(w, "id\t%s\n", result.ID)
	fmt.Fprintf(w, "name\t%s\n", result.Name)
	if result.Description != nil {
		fmt.Fprintf(w, "description\t%s\n", *result.Description)
	} else {
		fmt.Fprintf(w, "description\t\n")
	}
	fmt.Fprintf(w, "createdAt\t%s\n", result.CreatedAt)
	fmt.Fprintf(w, "updatedAt\t%s\n", result.UpdatedAt)
	return w.Flush()
}
