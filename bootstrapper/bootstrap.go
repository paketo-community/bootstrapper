package bootstrapper

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

func Bootstrap() error {
	var config map[string]string
	configFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}

	name := config["buildpack"]
	templatedPath := filepath.Join("/tmp", name)

	err = copyTemplateToTempDir(templatedPath)

	if err != nil {
		return err
	}

	err = filepath.Walk(templatedPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if strings.HasPrefix(path, filepath.Join(templatedPath, "bin")) ||
			strings.HasPrefix(path, filepath.Join(templatedPath, "vendor")) ||
			strings.HasPrefix(path, filepath.Join(templatedPath, ".github")) ||
			strings.HasPrefix(path, filepath.Join(templatedPath, ".bin")) {
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

func copyTemplateToTempDir(path string) error {
	fmt.Println(path)
	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		return err
	}

	err = CopyDirectory("template-cnb", path)
	if err != nil {
		return err
	}

	return nil
}
