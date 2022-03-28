package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var bootstrapperVersion string

func version() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "version of bootstrapper",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("bootstrapper %s\n", bootstrapperVersion)
		},
	}

	return cmd
}

func init() {
	rootCmd.AddCommand(version())
}
