package bootstrapper

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/paketo-buildpacks/packit/fs"
	"gopkg.in/yaml.v2"
)

//go:generate faux --interface TemplateWriter --output fakes/template_writer.go
type TemplateWriter interface {
	FillOutTemplate(path string, config Config) error
}

type Config struct {
	Buildpack    string `yaml:"buildpack"`
	Organization string `yaml:"organization"`
}

func Bootstrap(templateWriter TemplateWriter, configPath, templatePath, outputPath string) error {
	var config Config

	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}

	if outputPath == "" {
		outputPath = filepath.Join("/tmp", config.Buildpack)
	}

	err = fs.Copy(templatePath, outputPath)
	if err != nil {
		panic(err)
	}

	err = filepath.Walk(outputPath, func(path string, info os.FileInfo, err error) error {
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
			panic(err)
			return err
		}

		return nil
	})

	return err
}
