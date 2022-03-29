package commands

import (
	"fmt"
	"os"

	"github.com/paketo-buildpacks/packit/v2/pexec"
	internal "github.com/paketo-community/bootstrapper/commands/internal"
	"github.com/spf13/cobra"
)

type flags struct {
	buildpackName string
	outputDir     string
	templateDir   string
}

func run() *cobra.Command {
	var flags flags
	cmd := &cobra.Command{
		Use:   "run",
		Short: "bootstrap a packit-compliant buildpack",
		RunE: func(cmd *cobra.Command, args []string) error {
			return bootstrapRun(flags)
		},
	}
	cmd.Flags().StringVar(&flags.buildpackName, "buildpack-name", "", "name of buildpack in the form of <organization>/<buildpack>")
	cmd.Flags().StringVar(&flags.outputDir, "output", "", "path to location of output CNB")
	cmd.Flags().StringVar(&flags.templateDir, "template", "template-cnb", "path to template CNB directory")

	err := cmd.MarkFlagRequired("buildpack-name")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to mark buildpack flag as required")
	}
	err = cmd.MarkFlagRequired("output")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to mark output flag as required")
	}
	return cmd
}

func init() {
	rootCmd.AddCommand(run())
}

func bootstrapRun(flags flags) error {
	templatizer := internal.NewTemplatizer()
	golang := pexec.NewExecutable("go")
	fmt.Printf("Bootstrapping %s buildpack from template at %s \n", flags.buildpackName, flags.templateDir)
	err := internal.Bootstrap(templatizer, flags.buildpackName, flags.outputDir, flags.templateDir, golang)
	if err != nil {
		return err
	}
	fmt.Printf("Success! Buildpack available at %s\n", flags.outputDir)
	return nil
}
