package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/gentij/taskforge/apps/cli/internal/api"
	"github.com/gentij/taskforge/apps/cli/internal/output"
	"github.com/spf13/cobra"
)

type workflowItem struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	IsActive        bool    `json:"isActive"`
	LatestVersionID *string `json:"latestVersionId"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

type paginationMeta struct {
	Page       int  `json:"page"`
	PageSize   int  `json:"pageSize"`
	Total      int  `json:"total"`
	TotalPages int  `json:"totalPages"`
	HasNext    bool `json:"hasNext"`
	HasPrev    bool `json:"hasPrev"`
}

type workflowListResponse struct {
	Items      []workflowItem `json:"items"`
	Pagination paginationMeta `json:"pagination"`
}

var workflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Manage workflows",
}

var workflowListPage int
var workflowListPageSize int
var workflowCreateName string
var workflowCreateDefinition string
var workflowUpdateName string
var workflowUpdateIsActive bool
var workflowRunInput string
var workflowRunOverrides string
var workflowValidateDefinition string

func init() {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List workflows",
		RunE:  workflowList,
	}
	listCmd.Flags().IntVar(&workflowListPage, "page", 1, "Page number")
	listCmd.Flags().IntVar(&workflowListPageSize, "page-size", 25, "Page size")

	getCmd := &cobra.Command{
		Use:   "get <workflow-id>",
		Short: "Get a workflow",
		Args:  cobra.ExactArgs(1),
		RunE:  workflowGet,
	}

	workflowCmd.AddCommand(listCmd)
	workflowCmd.AddCommand(getCmd)
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a workflow",
		RunE:  workflowCreate,
	}
	createCmd.Flags().StringVar(&workflowCreateName, "name", "", "Workflow name")
	createCmd.Flags().StringVar(&workflowCreateDefinition, "definition", "", "Path to definition JSON")
	_ = createCmd.MarkFlagRequired("name")
	_ = createCmd.MarkFlagRequired("definition")

	updateCmd := &cobra.Command{
		Use:   "update <workflow-id>",
		Short: "Update a workflow",
		Args:  cobra.ExactArgs(1),
		RunE:  workflowUpdate,
	}
	updateCmd.Flags().StringVar(&workflowUpdateName, "name", "", "Workflow name")
	updateCmd.Flags().BoolVar(&workflowUpdateIsActive, "is-active", false, "Set workflow active state")

	deleteCmd := &cobra.Command{
		Use:   "delete <workflow-id>",
		Short: "Delete a workflow (soft)",
		Args:  cobra.ExactArgs(1),
		RunE:  workflowDelete,
	}

	runCmd := &cobra.Command{
		Use:   "run <workflow-id>",
		Short: "Run a workflow",
		Args:  cobra.ExactArgs(1),
		RunE:  workflowRun,
	}
	runCmd.Flags().StringVar(&workflowRunInput, "input", "", "Path to input JSON")
	runCmd.Flags().StringVar(&workflowRunOverrides, "overrides", "", "Path to overrides JSON")

	workflowCmd.AddCommand(createCmd)
	workflowCmd.AddCommand(updateCmd)
	workflowCmd.AddCommand(deleteCmd)
	workflowCmd.AddCommand(runCmd)
	workflowCmd.AddCommand(workflowVersionCmd)

	validateCmd := &cobra.Command{
		Use:   "validate <workflow-id>",
		Short: "Validate a workflow definition",
		Args:  cobra.ExactArgs(1),
		RunE:  workflowValidate,
	}
	validateCmd.Flags().StringVar(&workflowValidateDefinition, "definition", "", "Path to definition JSON")
	_ = validateCmd.MarkFlagRequired("definition")
	workflowCmd.AddCommand(validateCmd)
}

func workflowList(cmd *cobra.Command, args []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.ServerURL, cfg.Token)
	query := url.Values{}
	query.Set("page", fmt.Sprintf("%d", workflowListPage))
	query.Set("pageSize", fmt.Sprintf("%d", workflowListPageSize))

	var result workflowListResponse
	if err := client.GetJSON("/workflows?"+query.Encode(), &result); err != nil {
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
	fmt.Fprintln(w, "ID\tNAME\tACTIVE\tLATEST_VERSION")
	for _, item := range result.Items {
		latest := ""
		if item.LatestVersionID != nil {
			latest = *item.LatestVersionID
		}
		fmt.Fprintf(w, "%s\t%s\t%t\t%s\n", item.ID, item.Name, item.IsActive, latest)
	}
	return w.Flush()
}

func workflowGet(cmd *cobra.Command, args []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.ServerURL, cfg.Token)
	var result workflowItem
	if err := client.GetJSON("/workflows/"+args[0], &result); err != nil {
		return err
	}

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
	fmt.Fprintf(w, "isActive\t%t\n", result.IsActive)
	if result.LatestVersionID != nil {
		fmt.Fprintf(w, "latestVersionId\t%s\n", *result.LatestVersionID)
	} else {
		fmt.Fprintf(w, "latestVersionId\t\n")
	}
	fmt.Fprintf(w, "createdAt\t%s\n", result.CreatedAt)
	fmt.Fprintf(w, "updatedAt\t%s\n", result.UpdatedAt)
	return w.Flush()
}

func workflowCreate(cmd *cobra.Command, args []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	definition, err := readJSONFile(workflowCreateDefinition)
	if err != nil {
		return err
	}

	payload := map[string]any{
		"name":       workflowCreateName,
		"definition": definition,
	}

	client := api.NewClient(cfg.ServerURL, cfg.Token)
	var result workflowItem
	if err := client.PostJSON("/workflows", payload, &result); err != nil {
		return err
	}

	return printWorkflow(result)
}

func workflowUpdate(cmd *cobra.Command, args []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	patch := map[string]any{}
	if workflowUpdateName != "" {
		patch["name"] = workflowUpdateName
	}
	if cmd.Flags().Changed("is-active") {
		patch["isActive"] = workflowUpdateIsActive
	}
	if len(patch) == 0 {
		return fmt.Errorf("no fields to update")
	}

	client := api.NewClient(cfg.ServerURL, cfg.Token)
	var result workflowItem
	if err := client.PatchJSON("/workflows/"+args[0], patch, &result); err != nil {
		return err
	}

	return printWorkflow(result)
}

func workflowDelete(cmd *cobra.Command, args []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.ServerURL, cfg.Token)
	var result workflowItem
	if err := client.DeleteJSON("/workflows/"+args[0], &result); err != nil {
		return err
	}

	return printWorkflow(result)
}

func workflowRun(cmd *cobra.Command, args []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	input, err := readOptionalJSONFile(workflowRunInput)
	if err != nil {
		return err
	}

	overrides, err := readOptionalJSONFile(workflowRunOverrides)
	if err != nil {
		return err
	}

	payload := map[string]any{
		"input":     input,
		"overrides": overrides,
	}

	client := api.NewClient(cfg.ServerURL, cfg.Token)
	var result struct {
		WorkflowRunID string `json:"workflowRunId"`
		Status        string `json:"status"`
	}
	if err := client.PostJSON("/workflows/"+args[0]+"/run", payload, &result); err != nil {
		return err
	}

	if outputJSON {
		return output.PrintJSON(result)
	}
	if quiet {
		fmt.Fprintln(os.Stdout, result.WorkflowRunID)
		return nil
	}

	w := output.NewTableWriter()
	fmt.Fprintln(w, "FIELD\tVALUE")
	fmt.Fprintf(w, "workflowRunId\t%s\n", result.WorkflowRunID)
	fmt.Fprintf(w, "status\t%s\n", result.Status)
	return w.Flush()
}

func workflowValidate(cmd *cobra.Command, args []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	definition, err := readJSONFile(workflowValidateDefinition)
	if err != nil {
		return err
	}

	payload := map[string]any{
		"definition": definition,
	}

	client := api.NewClient(cfg.ServerURL, cfg.Token)
	var result struct {
		Valid                bool                `json:"valid"`
		Issues               []any               `json:"issues"`
		InferredDependencies map[string][]string `json:"inferredDependencies"`
		ExecutionBatches     [][]string          `json:"executionBatches"`
		ReferencedSecrets    []string            `json:"referencedSecrets"`
	}
	if err := client.PostJSON("/workflows/"+args[0]+"/versions/validate", payload, &result); err != nil {
		return err
	}

	if outputJSON {
		return output.PrintJSON(result)
	}

	if quiet {
		fmt.Fprintln(os.Stdout, result.Valid)
		return nil
	}

	issueCount := 0
	if result.Issues != nil {
		issueCount = len(result.Issues)
	}

	w := output.NewTableWriter()
	fmt.Fprintln(w, "FIELD\tVALUE")
	fmt.Fprintf(w, "valid\t%t\n", result.Valid)
	fmt.Fprintf(w, "issues\t%d\n", issueCount)
	return w.Flush()
}

func printWorkflow(result workflowItem) error {
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
	fmt.Fprintf(w, "isActive\t%t\n", result.IsActive)
	if result.LatestVersionID != nil {
		fmt.Fprintf(w, "latestVersionId\t%s\n", *result.LatestVersionID)
	} else {
		fmt.Fprintf(w, "latestVersionId\t\n")
	}
	fmt.Fprintf(w, "createdAt\t%s\n", result.CreatedAt)
	fmt.Fprintf(w, "updatedAt\t%s\n", result.UpdatedAt)
	return w.Flush()
}

func readJSONFile(path string) (any, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		return nil, err
	}

	return value, nil
}

func readOptionalJSONFile(path string) (any, error) {
	if strings.TrimSpace(path) == "" {
		return map[string]any{}, nil
	}

	return readJSONFile(path)
}
