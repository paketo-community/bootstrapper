package template

import (
	"github.com/cloudfoundry/packit"
)

func Build() packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		return packit.BuildResult{}, nil
	}
}
