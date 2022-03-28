package commands

import "github.com/spf13/cobra"

var (
	rootCmd = &cobra.Command{
		Use: "bootstrapper",
	}
)

func Execute() error {
	return rootCmd.Execute()
}
