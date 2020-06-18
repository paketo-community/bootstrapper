package main

import (
	"github.com/paketo-buildpacks/packit"
	"github.com/{{ .organization }}/{{ .buildpack }}"
)

func main() {
	packit.Run({{ .buildpack }}.Detect(), {{ .buildpack }}.Build())
}
