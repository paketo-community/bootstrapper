package main_test

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

var path string

func TestBootstrapperTool(t *testing.T) {
	// Need the 30 second timeout to run tool, run integration tests and unit tests in created buildpack
	SetDefaultEventuallyTimeout(30 * time.Second)

	suite := spec.New("Bootstrapper", spec.Report(report.Terminal{}))
	suite.Before(func(t *testing.T) {
		var (
			Expect = NewWithT(t).Expect
			err    error
		)

		path, err = gexec.Build("github.com/paketo-community/bootstrapper", "-ldflags", `-X github.com/paketo-community/bootstrapper/commands.bootstrapperVersion=1.2.3`)
		Expect(err).NotTo(HaveOccurred())
	})
	suite("Run", testRun)
	suite("Version", testVersion)
	suite.Run(t)
}
