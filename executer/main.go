package main

import (
	"flag"
	"fmt"

	"github.com/paketo-community/bootstrapper/bootstrapper"
)

func main() {
	var (
		buildpack    string
		templatePath string
		outputPath   string
	)

	flag.StringVar(&buildpack, "buildpack", "", "<org>/<name>")
	flag.StringVar(&templatePath, "template-path", "template-cnb", "path to the template")
	flag.StringVar(&outputPath, "output-path", "", "path to the new cnb")
	flag.Parse()

	templatizer := bootstrapper.NewTemplatizer()
	err := bootstrapper.Bootstrap(templatizer, buildpack, templatePath, outputPath)

	if err != nil {
		fmt.Println(err)
	}
}
