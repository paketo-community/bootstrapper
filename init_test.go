package main_test

import (
	"testing"
	"time"

	"github.com/onsi/gomega/gexec"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
)

var path string

func TestBootstrapperTool(t *testing.T) {
	SetDefaultEventuallyTimeout(60 * time.Second)

	var (
		err    error
		Expect = NewWithT(t).Expect
	)
	path, err = gexec.Build("github.com/paketo-community/bootstrapper", "-ldflags", `-X github.com/paketo-community/bootstrapper/commands.bootstrapperVersion=1.2.3`)
	Expect(err).NotTo(HaveOccurred())

	suite := spec.New("Bootstrapper", spec.Report(report.Terminal{}))
	suite("Run", testRun)
	suite("Version", testVersion)
	suite.Run(t)
}
