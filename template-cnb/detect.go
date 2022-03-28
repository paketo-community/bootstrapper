package {{ .Buildpack | RemoveHyphens}}

import "github.com/paketo-buildpacks/packit/v2"

func Detect() packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		return packit.DetectResult{}, nil
	}
}
