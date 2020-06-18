package bootstrapper

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

func Bootstrap() error {
	var (
		config       map[string]string
		configPath   string
		templatePath string
		outputPath   string
	)

	flag.StringVar(&configPath, "config-path", "config.yml", "path to the config file")
	flag.StringVar(&templatePath, "template-path", "template-cnb", "path to the template")
	flag.StringVar(&outputPath, "output-path", "", "path to the new cnb")
	flag.Parse()

	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}

	name := config["buildpack"]

	if outputPath == "" {
		outputPath = filepath.Join("/tmp", name)
	}

	err = copyTemplateToTempDir(outputPath, templatePath)

	if err != nil {
		panic(err)
	}

	err = filepath.Walk(outputPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if strings.HasPrefix(path, filepath.Join(outputPath, "bin")) ||
			strings.HasPrefix(path, filepath.Join(outputPath, "vendor")) ||
			strings.HasPrefix(path, filepath.Join(outputPath, ".github")) ||
			strings.HasPrefix(path, filepath.Join(outputPath, ".bin")) {
			return nil
		}

		buildpackTOML, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}

		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			panic(err)
		}

		funcMap := template.FuncMap{
			"Title": strings.Title,
		}
		t := template.Must(template.New("t1").Funcs(funcMap).Parse(string(buildpackTOML)))

		err = t.Execute(file, config)
		if err != nil {
			panic(err)
		}

		file.Close()
		return nil

	})
	return err
}

func copyTemplateToTempDir(path, templatePath string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	err = CopyDirectory(templatePath, path)
	if err != nil {
		return err
	}

	return nil
}
