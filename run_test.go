package main_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testRun(t *testing.T, context spec.G, it spec.S) {
	var (
		withT        = NewWithT(t)
		Expect       = withT.Expect
		Eventually   = withT.Eventually
		outputPath   string
		templatePath string
		buffer       *bytes.Buffer
	)

	it.Before(func() {
		var err error
		outputPath, err = os.MkdirTemp("", "")
		Expect(err).NotTo(HaveOccurred())

		templatePath, err = os.MkdirTemp("", "")
		Expect(err).NotTo(HaveOccurred())

		err = os.WriteFile(filepath.Join(templatePath, "go.mod"), []byte("module github.com/test/test"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		buffer = bytes.NewBuffer(nil)
	})

	it.After(func() {
		Expect(os.RemoveAll(outputPath)).To(Succeed())
		Expect(os.RemoveAll(templatePath)).To(Succeed())
	})

	context("bootstrapping a buildpack from default template", func() {
		it("creates a buildpack with runnable tests and scripts in the output directory", func() {
			command := exec.Command(
				path, "run",
				"--buildpack-name", "some-org/some-buildpack",
				"--output", outputPath,
			)
			session, err := gexec.Start(command, buffer, buffer)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0), func() string { return buffer.String() })

			Expect(session.Out).To(gbytes.Say("Bootstrapping some-org/some-buildpack buildpack from template at template-cnb"))
			Expect(session.Out).To(gbytes.Say(fmt.Sprintf("Success. Buildpack available at %s", outputPath)))

			Expect(outputPath).To(BeADirectory())
			b, err := os.ReadFile(filepath.Join(outputPath, "go.mod"))
			Expect(err).NotTo(HaveOccurred())
			Eventually(string(b)).Should(ContainSubstring("module github.com/some-org/some-buildpack"))

			unitCmd := exec.Command(filepath.Join(outputPath, "scripts", "unit.sh"))
			unitCmd.Dir = outputPath
			unitBuffer := gbytes.NewBuffer()

			unitSession, err := gexec.Start(unitCmd, unitBuffer, unitBuffer)
			Expect(err).NotTo(HaveOccurred())

			Eventually(unitSession).Should(gexec.Exit(0), func() string { return fmt.Sprintf("output:\n%s\n", unitBuffer.Contents()) })

			integrationCmd := exec.Command(filepath.Join(outputPath, "scripts", "integration.sh"))
			integrationCmd.Dir = outputPath
			integrationBuffer := gbytes.NewBuffer()

			integrationSession, err := gexec.Start(integrationCmd, integrationBuffer, integrationBuffer)
			Expect(err).NotTo(HaveOccurred())

			Eventually(integrationSession).Should(gexec.Exit(1))
			Expect(string(integrationBuffer.Contents())).To(ContainSubstring("Not Implemented"))
		})
	})

	context("with a custom template that contains a go.mod", func() {
		it("creates a custom-template based buildpack in the output directory", func() {
			command := exec.Command(
				path, "run",
				"--buildpack-name", "some-org/someBuildpack",
				"--output", outputPath,
				"--template", templatePath,
			)
			session, err := gexec.Start(command, buffer, buffer)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0), func() string { return buffer.String() })

			Expect(session.Out).To(gbytes.Say(fmt.Sprintf("Bootstrapping some-org/someBuildpack buildpack from template at %s", templatePath)))
			Expect(session.Out).To(gbytes.Say(fmt.Sprintf("Success. Buildpack available at %s", outputPath)))

			Expect(outputPath).To(BeADirectory())
			Eventually(filepath.Join(outputPath, "go.mod")).Should(BeAnExistingFile())
			b, err := os.ReadFile(filepath.Join(outputPath, "go.mod"))
			Expect(err).NotTo(HaveOccurred())
			Eventually(string(b)).Should(ContainSubstring("module github.com/some-org/someBuildpack"))
		})
	})

	context("failure cases", func() {
		context("when the all the required flags are not set", func() {
			it("prints an error message", func() {
				command := exec.Command(path, "run")
				session, err := gexec.Start(command, buffer, buffer)
				Expect(err).NotTo(HaveOccurred())
				Eventually(session).Should(gexec.Exit(1), func() string { return buffer.String() })
				Expect(string(session.Err.Contents())).To(ContainSubstring("failed to execute: required flag(s) \"buildpack-name\", \"output\" not set"))
			})
		})

		context("when the required buildpack-name flag is not set", func() {
			it("prints an error message", func() {
				command := exec.Command(
					path, "run",
					"--output", outputPath,
					"--template", templatePath,
				)
				session, err := gexec.Start(command, buffer, buffer)
				Expect(err).NotTo(HaveOccurred())
				Eventually(session).Should(gexec.Exit(1), func() string { return buffer.String() })
				Expect(string(session.Err.Contents())).To(ContainSubstring("failed to execute: required flag(s) \"buildpack-name\" not set"))
			})
		})

		context("when the required output flag is not set", func() {
			it("prints an error message", func() {
				command := exec.Command(
					path, "run",
					"--buildpack-name", "some-org/someBuildpack",
					"--template", templatePath,
				)
				session, err := gexec.Start(command, buffer, buffer)
				Expect(err).NotTo(HaveOccurred())
				Eventually(session).Should(gexec.Exit(1), func() string { return buffer.String() })
				Expect(string(session.Err.Contents())).To(ContainSubstring("failed to execute: required flag(s) \"output\" not set"))
			})
		})
	})
}
