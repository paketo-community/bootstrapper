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
	bootstrapper, err = gexec.Build("github.com/joshzarrabi/cnb-bootstrapper/executer")
	Expect(err).NotTo(HaveOccurred())
	SetDefaultEventuallyTimeout(5 * time.Second)

	spec.Run(t, "dispatch", func(t *testing.T, context spec.G, it spec.S) {
		var (
			Expect     = NewWithT(t).Expect
			Eventually = NewWithT(t).Eventually
		)

		context("when given a valid config", func() {
			var (
				configPath string
				outputPath string
			)

			it.Before(func() {
				outputPath, err = ioutil.TempDir("", "")
				Expect(err).NotTo(HaveOccurred())

				configFile, err := ioutil.TempFile("", "config.yml")
				Expect(err).NotTo(HaveOccurred())

				_, err = configFile.WriteString(`---
organization: some-org
buildpack: someBuildpack
`)
				Expect(err).NotTo(HaveOccurred())

				configPath = configFile.Name()
			})

			it.After(func() {
				Expect(os.RemoveAll(outputPath)).To(Succeed())
				Expect(os.RemoveAll(configPath)).To(Succeed())
			})

			it("creates a buildpack that can run `./scripts/package`", func() {
				command := exec.Command(
					bootstrapper,
					"--config-path", configPath,
					"--template-path", "/home/arjun/workspace/cnb-bootstrapper/template-cnb",
					"--output-path", outputPath,
				)
				buffer := gbytes.NewBuffer()

				session, err := gexec.Start(command, buffer, buffer)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(0), func() string { return fmt.Sprintf("output:\n%s\n", buffer.Contents()) })

				packageCmd := exec.Command(filepath.Join("scripts", "package.sh"))
				packageCmd.Dir = outputPath
				packageBuffer := gbytes.NewBuffer()

				packageSession, err := gexec.Start(packageCmd, packageBuffer, packageBuffer)
				Expect(err).NotTo(HaveOccurred())

				Eventually(packageSession).Should(gexec.Exit(0), func() string { return fmt.Sprintf("output:\n%s\n", packageBuffer.Contents()) })
			})
		})
	}, spec.Report(report.Terminal{}), spec.Parallel())
}
