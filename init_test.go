package bootstrapper_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestBootstrapper(t *testing.T) {
	suite := spec.New("bootstrapper", spec.Report(report.Terminal{}))
	suite("Bootstrap", testBootstrap)
	suite("Templatizer", testTemplatizer)
	suite.Run(t)
}
