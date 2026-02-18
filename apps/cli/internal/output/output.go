package output

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/gentij/taskforge/apps/cli/internal/api"
)

type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
)

func PrintJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(os.Stdout, string(data))
	return err
}

func NewTableWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 2, 4, 2, ' ', 0)
}

func BoolLabel(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

func PrintPagination(meta api.Pagination) error {
	if meta.TotalPages <= 1 {
		return nil
	}

	_, err := fmt.Fprintf(os.Stdout, "Page %d/%d · Total %d\n", meta.Page, meta.TotalPages, meta.Total)
	return err
}

func PrintError(err error) {
	if err == nil {
		return
	}

	if apiErr := api.AsAPIError(err); apiErr != nil {
		fmt.Fprintf(os.Stderr, "ERROR %s: %s\n", apiErr.Code, apiErr.Message)
		if apiErr.Details != nil {
			data, marshalErr := json.MarshalIndent(apiErr.Details, "", "  ")
			if marshalErr == nil {
				fmt.Fprintln(os.Stderr, string(data))
			}
		}
		return
	}

	fmt.Fprintf(os.Stderr, "ERROR %s\n", err.Error())
}
