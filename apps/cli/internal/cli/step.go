package cli

import (
	"fmt"
	"os"

	"github.com/gentij/taskforge/apps/cli/internal/output"
	"github.com/spf13/cobra"
)

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
	ctx := GetContext(cmd.Context())
	if ctx == nil {
		return fmt.Errorf("missing context")
	}

	workflowID := args[0]
	runID := args[1]
	result, err := ctx.Client.ListStepRuns(workflowID, runID, stepListPage, stepListPageSize)
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
		started := ""
		if item.StartedAt != nil {
			started = *item.StartedAt
		}
		rows = append(rows, []string{item.ID, item.StepKey, item.Status, started})
	}
	return output.PrintListTable([]string{"ID", "STEP_KEY", "STATUS", "STARTED"}, rows)
}

func stepGet(cmd *cobra.Command, args []string) error {
	ctx := GetContext(cmd.Context())
	if ctx == nil {
		return fmt.Errorf("missing context")
	}

	workflowID := args[0]
	runID := args[1]
	stepID := args[2]
	result, err := ctx.Client.GetStepRun(workflowID, runID, stepID)
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
		{"workflowRunId", result.WorkflowRunID},
		{"stepKey", result.StepKey},
		{"status", result.Status},
		{"startedAt", started},
		{"finishedAt", finished},
	})
}
