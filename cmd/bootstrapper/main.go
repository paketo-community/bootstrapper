package main

import (
	"flag"
	"log"

	"github.com/paketo-buildpacks/packit/pexec"
	"github.com/paketo-community/bootstrapper"
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
	golang := pexec.NewExecutable("go")

	err := bootstrapper.Bootstrap(templatizer, buildpack, templatePath, outputPath, golang)
	if err != nil {
		log.Fatalln(err)
	}
}
