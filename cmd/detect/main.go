package main

import (
	"{{ .buildpack }}/{{ .buildpack }}"

	"github.com/cloudfoundry/packit"
)

func main() {
	packit.Detect({{ .buildpack }}.Detect())
}
