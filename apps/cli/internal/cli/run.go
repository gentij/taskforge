package cli

import (
	"fmt"
	"net/url"
	"os"

	"github.com/gentij/taskforge/apps/cli/internal/api"
	"github.com/gentij/taskforge/apps/cli/internal/output"
	"github.com/spf13/cobra"
)

type runItem struct {
	ID                string  `json:"id"`
	WorkflowID        string  `json:"workflowId"`
	WorkflowVersionID string  `json:"workflowVersionId"`
	Status            string  `json:"status"`
	CreatedAt         string  `json:"createdAt"`
	StartedAt         *string `json:"startedAt"`
	FinishedAt        *string `json:"finishedAt"`
}

type runListResponse struct {
	Items      []runItem      `json:"items"`
	Pagination paginationMeta `json:"pagination"`
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Manage workflow runs",
}

var runListPage int
var runListPageSize int

func init() {
	listCmd := &cobra.Command{
		Use:   "list <workflow-id>",
		Short: "List workflow runs",
		Args:  cobra.ExactArgs(1),
		RunE:  runList,
	}
	listCmd.Flags().IntVar(&runListPage, "page", 1, "Page number")
	listCmd.Flags().IntVar(&runListPageSize, "page-size", 25, "Page size")

	getCmd := &cobra.Command{
		Use:   "get <workflow-id> <run-id>",
		Short: "Get a workflow run",
		Args:  cobra.ExactArgs(2),
		RunE:  runGet,
	}

	runCmd.AddCommand(listCmd)
	runCmd.AddCommand(getCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	workflowID := args[0]
	client := api.NewClient(cfg.ServerURL, cfg.Token)
	query := url.Values{}
	query.Set("page", fmt.Sprintf("%d", runListPage))
	query.Set("pageSize", fmt.Sprintf("%d", runListPageSize))

	var result runListResponse
	if err := client.GetJSON("/workflows/"+workflowID+"/runs?"+query.Encode(), &result); err != nil {
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
	fmt.Fprintln(w, "ID\tSTATUS\tVERSION\tCREATED")
	for _, item := range result.Items {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", item.ID, item.Status, item.WorkflowVersionID, item.CreatedAt)
	}
	return w.Flush()
}

func runGet(cmd *cobra.Command, args []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	workflowID := args[0]
	runID := args[1]
	client := api.NewClient(cfg.ServerURL, cfg.Token)
	var result runItem
	if err := client.GetJSON("/workflows/"+workflowID+"/runs/"+runID, &result); err != nil {
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
	fmt.Fprintf(w, "workflowId\t%s\n", result.WorkflowID)
	fmt.Fprintf(w, "workflowVersionId\t%s\n", result.WorkflowVersionID)
	fmt.Fprintf(w, "status\t%s\n", result.Status)
	fmt.Fprintf(w, "createdAt\t%s\n", result.CreatedAt)
	if result.StartedAt != nil {
		fmt.Fprintf(w, "startedAt\t%s\n", *result.StartedAt)
	} else {
		fmt.Fprintf(w, "startedAt\t\n")
	}
	if result.FinishedAt != nil {
		fmt.Fprintf(w, "finishedAt\t%s\n", *result.FinishedAt)
	} else {
		fmt.Fprintf(w, "finishedAt\t\n")
	}
	return w.Flush()
}
