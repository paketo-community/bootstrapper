package internal_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestInternal(t *testing.T) {
	suite := spec.New("TestInternal", spec.Report(report.Terminal{}))
	suite("Bootstrap", testBootstrap)
	suite("Templatizer", testTemplatizer)
	suite.Run(t)
}
