package bootstrapper

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/paketo-buildpacks/packit/fs"
	"github.com/paketo-buildpacks/packit/pexec"
)

//go:generate faux --interface TemplateWriter --output fakes/template_writer.go
type TemplateWriter interface {
	FillOutTemplate(path string, config Config) error
}

//go:generate faux --interface Executable --output fakes/executable.go
type Executable interface {
	Execute(pexec.Execution) error
}

type Config struct {
	Organization string `yaml:"organization"`
	Buildpack    string `yaml:"buildpack"`
}

func Bootstrap(templateWriter TemplateWriter, buildpack, templatePath, outputPath string, golang Executable) error {
	parts := strings.Split(buildpack, "/")
	if len(parts) != 2 {
		return errors.New("buildpack name must be in format <organization>/<buildpack-name>")
	}

	config := Config{
		Organization: parts[0],
		Buildpack:    parts[1],
	}

	if outputPath == "" {
		outputPath = filepath.Join("/tmp", config.Buildpack)
	}

	err := fs.Copy(templatePath, outputPath)
	if err != nil {
		return fmt.Errorf("failed to copy template to the output path: %q", err)
	}

	err = filepath.Walk(outputPath, func(path string, info os.FileInfo, err error) error {
		switch {
		case err != nil:
			return err

		// NOTE: do nothing in these cases
		case info.IsDir():
		case strings.HasPrefix(path, filepath.Join(outputPath, ".github")):
		case strings.HasPrefix(path, filepath.Join(outputPath, "scripts")):

		default:
			err = templateWriter.FillOutTemplate(path, config)
			if err != nil {
				return fmt.Errorf("failed to fill out template: %q", err)
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	err = golang.Execute(pexec.Execution{
		Args: []string{"mod", "tidy"},
		Dir:  outputPath,
	})
	if err != nil {
		return fmt.Errorf("failed to run 'go mod tidy': %s", err)
	}

	return nil
}
