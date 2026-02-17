package output

import (
	"fmt"
)

func PrintListTable(headers []string, rows [][]string) error {
	w := NewTableWriter()

	for i, header := range headers {
		if i == len(headers)-1 {
			fmt.Fprintln(w, header)
		} else {
			fmt.Fprint(w, header+"\t")
		}
	}

	for _, row := range rows {
		for i, col := range row {
			if i == len(row)-1 {
				fmt.Fprintln(w, col)
			} else {
				fmt.Fprint(w, col+"\t")
			}
		}
	}

	return w.Flush()
}

func PrintKVTable(pairs [][2]string) error {
	w := NewTableWriter()
	fmt.Fprintln(w, "FIELD\tVALUE")
	for _, pair := range pairs {
		fmt.Fprintf(w, "%s\t%s\n", pair[0], pair[1])
	}
	return w.Flush()
}
