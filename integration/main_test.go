package main_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
)

var bootstrapper string

func TestBootstrapper(t *testing.T) {
	var Expect = NewWithT(t).Expect
	var err error
	bootstrapper, err = gexec.Build("github.com/paketo-community/bootstrapper/executer")
	Expect(err).NotTo(HaveOccurred())
	SetDefaultEventuallyTimeout(60 * time.Second)

	spec.Run(t, "dispatch", func(t *testing.T, context spec.G, it spec.S) {
		var (
			Expect     = NewWithT(t).Expect
			Eventually = NewWithT(t).Eventually
		)

		context("when given a valid config", func() {
			var outputPath string

			it.Before(func() {
				outputPath, err = ioutil.TempDir("", "")
				Expect(err).NotTo(HaveOccurred())
			})

			it.After(func() {
				Expect(os.RemoveAll(outputPath)).To(Succeed())
			})

			it("creates a buildpack that can run `./scripts/unit.sh` and `./scripts/integration.sh`", func() {
				command := exec.Command(
					bootstrapper,
					"--buildpack", "some-org/someBuildpack",
					"--template-path", "../template-cnb",
					"--output-path", outputPath,
				)
				buffer := gbytes.NewBuffer()

				session, err := gexec.Start(command, buffer, buffer)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(0), func() string { return fmt.Sprintf("output:\n%s\n", buffer.Contents()) })

				unitCmd := exec.Command(filepath.Join("scripts", "unit.sh"))
				unitCmd.Dir = outputPath
				unitBuffer := gbytes.NewBuffer()

				unitSession, err := gexec.Start(unitCmd, unitBuffer, unitBuffer)
				Expect(err).NotTo(HaveOccurred())

				Eventually(unitSession).Should(gexec.Exit(0), func() string { return fmt.Sprintf("output:\n%s\n", unitBuffer.Contents()) })

				integrationCmd := exec.Command(filepath.Join("scripts", "integration.sh"))
				integrationCmd.Dir = outputPath
				integrationBuffer := gbytes.NewBuffer()

				integrationSession, err := gexec.Start(integrationCmd, integrationBuffer, integrationBuffer)
				Expect(err).NotTo(HaveOccurred())

				Eventually(integrationSession).Should(gexec.Exit(1))
				Expect(integrationBuffer.Contents()).To(ContainSubstring("Not Implemented"))
			})
		})

		context("when buildpack name contains a hyphen", func() {
			var outputPath string

			it.Before(func() {
				outputPath, err = ioutil.TempDir("", "")
				Expect(err).NotTo(HaveOccurred())
			})

			it.After(func() {
				Expect(os.RemoveAll(outputPath)).To(Succeed())
			})

			it("creates a buildpack that can run `./scripts/unit.sh` and `./scripts/integration.sh`", func() {
				command := exec.Command(
					bootstrapper,
					"--buildpack", "some-org/some-hyphenated-buildpack",
					"--template-path", "../template-cnb",
					"--output-path", outputPath,
				)
				buffer := gbytes.NewBuffer()

				session, err := gexec.Start(command, buffer, buffer)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(0), func() string { return fmt.Sprintf("output:\n%s\n", buffer.Contents()) })

				unitCmd := exec.Command(filepath.Join("scripts", "unit.sh"))
				unitCmd.Dir = outputPath
				unitBuffer := gbytes.NewBuffer()

				unitSession, err := gexec.Start(unitCmd, unitBuffer, unitBuffer)
				Expect(err).NotTo(HaveOccurred())

				Eventually(unitSession).Should(gexec.Exit(0), func() string { return fmt.Sprintf("output:\n%s\n", unitBuffer.Contents()) })

				integrationCmd := exec.Command(filepath.Join("scripts", "integration.sh"))
				integrationCmd.Dir = outputPath
				integrationBuffer := gbytes.NewBuffer()

				integrationSession, err := gexec.Start(integrationCmd, integrationBuffer, integrationBuffer)
				Expect(err).NotTo(HaveOccurred())

				Eventually(integrationSession).Should(gexec.Exit(1))
				Expect(integrationBuffer.Contents()).To(ContainSubstring("Not Implemented"))
			})
		})
	}, spec.Report(report.Terminal{}))
}
