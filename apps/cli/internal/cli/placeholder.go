package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newNotImplementedCmd(use string, short string) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("not implemented")
		},
	}
}
