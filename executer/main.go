package main

import (
	"flag"
	"fmt"

	"github.com/joshzarrabi/cnb-bootstrapper/bootstrapper"
)

func main() {
	var (
		configPath   string
		templatePath string
		outputPath   string
	)

	flag.StringVar(&configPath, "config-path", "config.yml", "path to the config file")
	flag.StringVar(&templatePath, "template-path", "template-cnb", "path to the template")
	flag.StringVar(&outputPath, "output-path", "", "path to the new cnb")
	flag.Parse()

	templatizer := bootstrapper.NewTemplatizer()
	err := bootstrapper.Bootstrap(templatizer, configPath, templatePath, outputPath)

	if err != nil {
		fmt.Println(err)
	}
}
