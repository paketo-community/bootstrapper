package main

import (
	"github.com/paketo-buildpacks/packit"
	"github.com/{{ .Organization }}/{{ .Buildpack }}"
)

func main() {
	packit.Run({{ .Buildpack }}.Detect(), {{ .Buildpack }}.Build())
}
