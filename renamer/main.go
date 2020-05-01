package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if strings.HasPrefix(path, ".") || strings.HasPrefix(path, "bin") || strings.HasPrefix(path, "vendor") {
			return nil
		}

		dict := map[string]string{
			"buildpack": "template",
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

		err = t.Execute(file, dict)
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
