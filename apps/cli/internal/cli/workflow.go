package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"

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

	workflowCmd.AddCommand(createCmd)
	workflowCmd.AddCommand(newNotImplementedCmd("update", "Update a workflow"))
	workflowCmd.AddCommand(newNotImplementedCmd("delete", "Delete a workflow"))
	workflowCmd.AddCommand(newNotImplementedCmd("run", "Run a workflow"))
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
