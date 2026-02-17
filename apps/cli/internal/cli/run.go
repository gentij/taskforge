package cli

import (
	"fmt"
	"os"

	"github.com/gentij/taskforge/apps/cli/internal/output"
	"github.com/spf13/cobra"
)

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
	ctx := GetContext(cmd.Context())
	if ctx == nil {
		return fmt.Errorf("missing context")
	}

	workflowID := args[0]
	result, err := ctx.Client.ListWorkflowRuns(workflowID, runListPage, runListPageSize)
	if err != nil {
		return err
	}

	if ctx.OutputJSON {
		return output.PrintJSON(result)
	}

	if ctx.Quiet {
		for _, item := range result.Items {
			fmt.Fprintln(os.Stdout, item.ID)
		}
		return nil
	}

	rows := make([][]string, 0, len(result.Items))
	for _, item := range result.Items {
		rows = append(rows, []string{item.ID, item.Status, item.WorkflowVersionID, item.CreatedAt})
	}
	return output.PrintListTable([]string{"ID", "STATUS", "VERSION", "CREATED"}, rows)
}

func runGet(cmd *cobra.Command, args []string) error {
	ctx := GetContext(cmd.Context())
	if ctx == nil {
		return fmt.Errorf("missing context")
	}

	workflowID := args[0]
	runID := args[1]
	result, err := ctx.Client.GetWorkflowRun(workflowID, runID)
	if err != nil {
		return err
	}

	if ctx.OutputJSON {
		return output.PrintJSON(result)
	}
	if ctx.Quiet {
		fmt.Fprintln(os.Stdout, result.ID)
		return nil
	}

	started := ""
	if result.StartedAt != nil {
		started = *result.StartedAt
	}
	finished := ""
	if result.FinishedAt != nil {
		finished = *result.FinishedAt
	}

	return output.PrintKVTable([][2]string{
		{"id", result.ID},
		{"workflowId", result.WorkflowID},
		{"workflowVersionId", result.WorkflowVersionID},
		{"status", result.Status},
		{"createdAt", result.CreatedAt},
		{"startedAt", started},
		{"finishedAt", finished},
	})
}
