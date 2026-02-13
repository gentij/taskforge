package output

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
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
