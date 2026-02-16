package cli

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/gentij/taskforge/apps/cli/internal/api"
	"github.com/gentij/taskforge/apps/cli/internal/output"
	"github.com/spf13/cobra"
)

type triggerItem struct {
	ID         string `json:"id"`
	WorkflowID string `json:"workflowId"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	IsActive   bool   `json:"isActive"`
	Config     any    `json:"config"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

type triggerListResponse struct {
	Items      []triggerItem  `json:"items"`
	Pagination paginationMeta `json:"pagination"`
}

var triggerCmd = &cobra.Command{
	Use:   "trigger",
	Short: "Manage triggers",
}

var triggerListPage int
var triggerListPageSize int
var triggerCreateType string
var triggerCreateName string
var triggerCreateIsActive bool
var triggerCreateConfig string
var triggerUpdateName string
var triggerUpdateIsActive bool
var triggerUpdateConfig string

func init() {
	listCmd := &cobra.Command{
		Use:   "list <workflow-id>",
		Short: "List triggers",
		Args:  cobra.ExactArgs(1),
		RunE:  triggerList,
	}
	listCmd.Flags().IntVar(&triggerListPage, "page", 1, "Page number")
	listCmd.Flags().IntVar(&triggerListPageSize, "page-size", 25, "Page size")

	getCmd := &cobra.Command{
		Use:   "get <workflow-id> <trigger-id>",
		Short: "Get a trigger",
		Args:  cobra.ExactArgs(2),
		RunE:  triggerGet,
	}

	createCmd := &cobra.Command{
		Use:   "create <workflow-id>",
		Short: "Create a trigger",
		Args:  cobra.ExactArgs(1),
		RunE:  triggerCreate,
	}
	createCmd.Flags().StringVar(&triggerCreateType, "type", "", "Trigger type (MANUAL|CRON|WEBHOOK)")
	createCmd.Flags().StringVar(&triggerCreateName, "name", "", "Trigger name")
	createCmd.Flags().BoolVar(&triggerCreateIsActive, "is-active", true, "Trigger active state")
	createCmd.Flags().StringVar(&triggerCreateConfig, "config", "", "Path to config JSON")
	_ = createCmd.MarkFlagRequired("type")

	updateCmd := &cobra.Command{
		Use:   "update <workflow-id> <trigger-id>",
		Short: "Update a trigger",
		Args:  cobra.ExactArgs(2),
		RunE:  triggerUpdate,
	}
	updateCmd.Flags().StringVar(&triggerUpdateName, "name", "", "Trigger name")
	updateCmd.Flags().BoolVar(&triggerUpdateIsActive, "is-active", false, "Set trigger active state")
	updateCmd.Flags().StringVar(&triggerUpdateConfig, "config", "", "Path to config JSON")

	deleteCmd := &cobra.Command{
		Use:   "delete <workflow-id> <trigger-id>",
		Short: "Delete a trigger (soft)",
		Args:  cobra.ExactArgs(2),
		RunE:  triggerDelete,
	}

	triggerCmd.AddCommand(listCmd)
	triggerCmd.AddCommand(getCmd)
	triggerCmd.AddCommand(createCmd)
	triggerCmd.AddCommand(updateCmd)
	triggerCmd.AddCommand(deleteCmd)
}

func triggerList(cmd *cobra.Command, args []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	workflowID := args[0]
	client := api.NewClient(cfg.ServerURL, cfg.Token)
	query := url.Values{}
	query.Set("page", fmt.Sprintf("%d", triggerListPage))
	query.Set("pageSize", fmt.Sprintf("%d", triggerListPageSize))

	var result triggerListResponse
	if err := client.GetJSON("/workflows/"+workflowID+"/triggers?"+query.Encode(), &result); err != nil {
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
	fmt.Fprintln(w, "ID\tTYPE\tNAME\tACTIVE")
	for _, item := range result.Items {
		fmt.Fprintf(w, "%s\t%s\t%s\t%t\n", item.ID, item.Type, item.Name, item.IsActive)
	}
	return w.Flush()
}

func triggerGet(cmd *cobra.Command, args []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	workflowID := args[0]
	triggerID := args[1]
	client := api.NewClient(cfg.ServerURL, cfg.Token)
	var result triggerItem
	if err := client.GetJSON("/workflows/"+workflowID+"/triggers/"+triggerID, &result); err != nil {
		return err
	}

	return printTrigger(result)
}

func triggerCreate(cmd *cobra.Command, args []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	workflowID := args[0]
	configValue, err := readOptionalJSONFile(triggerCreateConfig)
	if err != nil {
		return err
	}

	payload := map[string]any{
		"type":     strings.ToUpper(triggerCreateType),
		"name":     triggerCreateName,
		"isActive": triggerCreateIsActive,
		"config":   configValue,
	}

	client := api.NewClient(cfg.ServerURL, cfg.Token)
	var result triggerItem
	if err := client.PostJSON("/workflows/"+workflowID+"/triggers", payload, &result); err != nil {
		return err
	}

	return printTrigger(result)
}

func triggerUpdate(cmd *cobra.Command, args []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	workflowID := args[0]
	triggerID := args[1]

	patch := map[string]any{}
	if triggerUpdateName != "" {
		patch["name"] = triggerUpdateName
	}
	if cmd.Flags().Changed("is-active") {
		patch["isActive"] = triggerUpdateIsActive
	}
	if strings.TrimSpace(triggerUpdateConfig) != "" {
		configValue, err := readJSONFile(triggerUpdateConfig)
		if err != nil {
			return err
		}
		patch["config"] = configValue
	}
	if len(patch) == 0 {
		return fmt.Errorf("no fields to update")
	}

	client := api.NewClient(cfg.ServerURL, cfg.Token)
	var result triggerItem
	if err := client.PatchJSON("/workflows/"+workflowID+"/triggers/"+triggerID, patch, &result); err != nil {
		return err
	}

	return printTrigger(result)
}

func triggerDelete(cmd *cobra.Command, args []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	workflowID := args[0]
	triggerID := args[1]
	client := api.NewClient(cfg.ServerURL, cfg.Token)
	var result triggerItem
	if err := client.DeleteJSON("/workflows/"+workflowID+"/triggers/"+triggerID, &result); err != nil {
		return err
	}

	return printTrigger(result)
}

func printTrigger(result triggerItem) error {
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
	fmt.Fprintf(w, "workflowId\t%s\n", result.WorkflowID)
	fmt.Fprintf(w, "type\t%s\n", result.Type)
	fmt.Fprintf(w, "name\t%s\n", result.Name)
	fmt.Fprintf(w, "isActive\t%t\n", result.IsActive)

	configData, err := json.Marshal(result.Config)
	if err == nil {
		fmt.Fprintf(w, "config\t%s\n", string(configData))
	} else {
		fmt.Fprintf(w, "config\t\n")
	}

	fmt.Fprintf(w, "createdAt\t%s\n", result.CreatedAt)
	fmt.Fprintf(w, "updatedAt\t%s\n", result.UpdatedAt)
	return w.Flush()
}
