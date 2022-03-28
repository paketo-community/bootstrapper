package main

import (
	"github.com/paketo-buildpacks/packit/v2"
	{{ .Buildpack | RemoveHyphens }} "github.com/{{ .Organization }}/{{ .Buildpack }}"
)

func main() {
	packit.Run({{ .Buildpack | RemoveHyphens }}.Detect(), {{ .Buildpack | RemoveHyphens }}.Build())
}
