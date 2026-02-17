package cli

import (
	"fmt"
	"net/url"
	"os"

	"github.com/gentij/taskforge/apps/cli/internal/api"
	"github.com/gentij/taskforge/apps/cli/internal/output"
	"github.com/spf13/cobra"
)

type workflowVersionItem struct {
	ID         string `json:"id"`
	WorkflowID string `json:"workflowId"`
	Version    int    `json:"version"`
	CreatedAt  string `json:"createdAt"`
}

type workflowVersionListResponse struct {
	Items      []workflowVersionItem `json:"items"`
	Pagination paginationMeta        `json:"pagination"`
}

var workflowVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Manage workflow versions",
}

var workflowVersionListPage int
var workflowVersionListPageSize int
var workflowVersionCreateDefinition string

func init() {
	listCmd := &cobra.Command{
		Use:   "list <workflow-id>",
		Short: "List workflow versions",
		Args:  cobra.ExactArgs(1),
		RunE:  workflowVersionList,
	}
	listCmd.Flags().IntVar(&workflowVersionListPage, "page", 1, "Page number")
	listCmd.Flags().IntVar(&workflowVersionListPageSize, "page-size", 25, "Page size")

	getCmd := &cobra.Command{
		Use:   "get <workflow-id> <version>",
		Short: "Get a workflow version",
		Args:  cobra.ExactArgs(2),
		RunE:  workflowVersionGet,
	}

	createCmd := &cobra.Command{
		Use:   "create <workflow-id>",
		Short: "Create a workflow version",
		Args:  cobra.ExactArgs(1),
		RunE:  workflowVersionCreate,
	}
	createCmd.Flags().StringVar(&workflowVersionCreateDefinition, "definition", "", "Path to definition JSON")
	_ = createCmd.MarkFlagRequired("definition")

	workflowVersionCmd.AddCommand(listCmd)
	workflowVersionCmd.AddCommand(getCmd)
	workflowVersionCmd.AddCommand(createCmd)
}

func workflowVersionList(cmd *cobra.Command, args []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	workflowID := args[0]
	client := api.NewClient(cfg.ServerURL, cfg.Token)
	query := url.Values{}
	query.Set("page", fmt.Sprintf("%d", workflowVersionListPage))
	query.Set("pageSize", fmt.Sprintf("%d", workflowVersionListPageSize))

	var result workflowVersionListResponse
	if err := client.GetJSON("/workflows/"+workflowID+"/versions?"+query.Encode(), &result); err != nil {
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
	fmt.Fprintln(w, "VERSION\tID\tCREATED")
	for _, item := range result.Items {
		fmt.Fprintf(w, "%d\t%s\t%s\n", item.Version, item.ID, item.CreatedAt)
	}
	return w.Flush()
}

func workflowVersionGet(cmd *cobra.Command, args []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	workflowID := args[0]
	version := args[1]
	client := api.NewClient(cfg.ServerURL, cfg.Token)
	var result workflowVersionItem
	if err := client.GetJSON("/workflows/"+workflowID+"/versions/"+version, &result); err != nil {
		return err
	}

	return printWorkflowVersion(result)
}

func workflowVersionCreate(cmd *cobra.Command, args []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	workflowID := args[0]
	definition, err := readJSONFile(workflowVersionCreateDefinition)
	if err != nil {
		return err
	}

	payload := map[string]any{
		"definition": definition,
	}

	client := api.NewClient(cfg.ServerURL, cfg.Token)
	var result workflowVersionItem
	if err := client.PostJSON("/workflows/"+workflowID+"/versions", payload, &result); err != nil {
		return err
	}

	return printWorkflowVersion(result)
}

func printWorkflowVersion(result workflowVersionItem) error {
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
	fmt.Fprintf(w, "version\t%d\n", result.Version)
	fmt.Fprintf(w, "createdAt\t%s\n", result.CreatedAt)
	return w.Flush()
}
