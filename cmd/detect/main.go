package main

import (
	"template/template"

	"github.com/cloudfoundry/packit"
)

func main() {
	packit.Detect(template.Detect())
}
