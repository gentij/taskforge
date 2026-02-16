package cli

import (
	"fmt"
	"net/url"
	"os"

	"github.com/gentij/taskforge/apps/cli/internal/api"
	"github.com/gentij/taskforge/apps/cli/internal/output"
	"github.com/spf13/cobra"
)

type stepRunItem struct {
	ID            string  `json:"id"`
	WorkflowRunID string  `json:"workflowRunId"`
	StepKey       string  `json:"stepKey"`
	Status        string  `json:"status"`
	CreatedAt     string  `json:"createdAt"`
	StartedAt     *string `json:"startedAt"`
	FinishedAt    *string `json:"finishedAt"`
}

type stepRunListResponse struct {
	Items      []stepRunItem  `json:"items"`
	Pagination paginationMeta `json:"pagination"`
}

var stepCmd = &cobra.Command{
	Use:   "step",
	Short: "Manage step runs",
}

var stepListPage int
var stepListPageSize int

func init() {
	listCmd := &cobra.Command{
		Use:   "list <workflow-id> <run-id>",
		Short: "List step runs",
		Args:  cobra.ExactArgs(2),
		RunE:  stepList,
	}
	listCmd.Flags().IntVar(&stepListPage, "page", 1, "Page number")
	listCmd.Flags().IntVar(&stepListPageSize, "page-size", 25, "Page size")

	getCmd := &cobra.Command{
		Use:   "get <workflow-id> <run-id> <step-run-id>",
		Short: "Get a step run",
		Args:  cobra.ExactArgs(3),
		RunE:  stepGet,
	}

	stepCmd.AddCommand(listCmd)
	stepCmd.AddCommand(getCmd)
}

func stepList(cmd *cobra.Command, args []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	workflowID := args[0]
	runID := args[1]
	client := api.NewClient(cfg.ServerURL, cfg.Token)
	query := url.Values{}
	query.Set("page", fmt.Sprintf("%d", stepListPage))
	query.Set("pageSize", fmt.Sprintf("%d", stepListPageSize))

	var result stepRunListResponse
	if err := client.GetJSON("/workflows/"+workflowID+"/runs/"+runID+"/steps?"+query.Encode(), &result); err != nil {
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
	fmt.Fprintln(w, "ID\tSTEP_KEY\tSTATUS\tSTARTED")
	for _, item := range result.Items {
		started := ""
		if item.StartedAt != nil {
			started = *item.StartedAt
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", item.ID, item.StepKey, item.Status, started)
	}
	return w.Flush()
}

func stepGet(cmd *cobra.Command, args []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	workflowID := args[0]
	runID := args[1]
	stepID := args[2]
	client := api.NewClient(cfg.ServerURL, cfg.Token)
	var result stepRunItem
	if err := client.GetJSON("/workflows/"+workflowID+"/runs/"+runID+"/steps/"+stepID, &result); err != nil {
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
	fmt.Fprintf(w, "workflowRunId\t%s\n", result.WorkflowRunID)
	fmt.Fprintf(w, "stepKey\t%s\n", result.StepKey)
	fmt.Fprintf(w, "status\t%s\n", result.Status)
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
