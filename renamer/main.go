package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

func main() {
	err := filepath.Walk("template-cnb", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if strings.HasPrefix(path, "template-cnb/.") || strings.HasPrefix(path, "template-cnb/bin") || strings.HasPrefix(path, "template-cnb/vendor") {
			return nil
		}

		var config map[string]string
		configFile, err := ioutil.ReadFile("config.yml")
		if err != nil {
			panic(err)
		}

		err = yaml.Unmarshal(configFile, &config)
		if err != nil {
			panic(err)
		}

		buildpackTOML, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}

		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			panic(err)
		}

		fmt.Println(path)

		t := template.Must(template.New("t1").Parse(string(buildpackTOML)))

		err = t.Execute(file, config)
		if err != nil {
			panic(err)
		}

		file.Close()
		return nil

	})

	if err != nil {
		panic(err)
	}
}
