package bootstrapper

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/paketo-buildpacks/packit/fs"
)

//go:generate faux --interface TemplateWriter --output fakes/template_writer.go
type TemplateWriter interface {
	FillOutTemplate(path string, config Config) error
}

type Config struct {
	Organization string `yaml:"organization"`
	Buildpack    string `yaml:"buildpack"`
}

func Bootstrap(templateWriter TemplateWriter, buildpack, templatePath, outputPath string) error {
	if len(strings.Split(buildpack, "/")) != 2 {
		return errors.New("buildpack name must be in format <organization>/<buildpack-name>")
	}

	config := Config{
		Organization: strings.Split(buildpack, "/")[0],
		Buildpack:    strings.Split(buildpack, "/")[1],
	}

	if outputPath == "" {
		outputPath = filepath.Join("/tmp", config.Buildpack)
	}

	err := fs.Copy(templatePath, outputPath)
	if err != nil {
		return fmt.Errorf("failed to copy template to the output path: %q", err)
	}

	err = filepath.Walk(outputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasPrefix(path, filepath.Join(outputPath, "bin")) ||
			strings.HasPrefix(path, filepath.Join(outputPath, ".github")) ||
			strings.HasPrefix(path, filepath.Join(outputPath, ".bin")) {
			return nil
		}

		err = templateWriter.FillOutTemplate(path, config)
		if err != nil {
			return fmt.Errorf("failed to fill out template: %q", err)
		}

		return nil
	})

	return err
}
