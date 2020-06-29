package {{ .Buildpack }}_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnit{{ .Buildpack | Title }}(t *testing.T) {
	suite := spec.New("{{ .Buildpack }}", spec.Report(report.Terminal{}), spec.Parallel())
	suite("Build", testBuild)
	suite("Detect", testDetect)
	suite.Run(t)
}
