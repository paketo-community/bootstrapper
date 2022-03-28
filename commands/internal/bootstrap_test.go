package internal_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/packit/v2/pexec"
	"github.com/paketo-community/bootstrapper/commands/internal"
	"github.com/paketo-community/bootstrapper/commands/internal/fakes"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testBootstrap(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		templatizer *fakes.TemplateWriter
		executable  *fakes.Executable

		outputPath   string
		templatePath string
		err          error

		templatizerCalls []struct {
			path   string
			config internal.Config
		}
	)

	it.Before(func() {
		outputPath, err = os.MkdirTemp("", "output-path")
		Expect(err).NotTo(HaveOccurred())

		templatePath, err = os.MkdirTemp("", "template-path")
		Expect(err).NotTo(HaveOccurred())

		err = os.WriteFile(filepath.Join(templatePath, "templ1"), []byte(""), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		err = os.WriteFile(filepath.Join(templatePath, "templ2"), []byte(""), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		templatizer = &fakes.TemplateWriter{}
		templatizer.FillOutTemplateCall.Stub = func(path string, config internal.Config) error {
			templatizerCalls = append(templatizerCalls, struct {
				path   string
				config internal.Config
			}{path: path, config: config})
			return nil
		}

		executable = &fakes.Executable{}
	})

	it.After(func() {
		Expect(os.RemoveAll(outputPath)).To(Succeed())
		Expect(os.RemoveAll(templatePath)).To(Succeed())
	})

	it("fills out every template in the correct path", func() {
		err := internal.Bootstrap(templatizer, "some-org/someBuildpack", outputPath, templatePath, executable)
		Expect(err).NotTo(HaveOccurred())

		Expect(templatizer.FillOutTemplateCall.CallCount).To(Equal(2))
		Expect(templatizerCalls[0].path).To(Equal(filepath.Join(outputPath, "templ1")))
		Expect(templatizerCalls[0].config).To(Equal(internal.Config{
			Buildpack:    "someBuildpack",
			Organization: "some-org",
		}))

		Expect(templatizerCalls[1].path).To(Equal(filepath.Join(outputPath, "templ2")))
		Expect(templatizerCalls[1].config).To(Equal(internal.Config{
			Buildpack:    "someBuildpack",
			Organization: "some-org",
		}))

		Expect(executable.ExecuteCall.Receives.Execution).To(Equal(pexec.Execution{
			Args: []string{"mod", "tidy"},
			Dir:  outputPath,
		}))
	})

	context("error cases", func() {
		context("when the buildpack name is malformed", func() {
			it("errors with a helpful message", func() {
				err := internal.Bootstrap(templatizer, "some-malformed-name", outputPath, templatePath, executable)
				Expect(err).To(MatchError(ContainSubstring("buildpack name must be in format <organization>/<buildpack-name>")))
			})
		})
		context("when the template can not be copied to the output path", func() {
			it.Before(func() {
				templatePath = "does-not-exist"
			})

			it("errors", func() {
				err := internal.Bootstrap(templatizer, "some-org/someBuildpack", outputPath, templatePath, executable)
				Expect(err).To(MatchError(ContainSubstring("failed to copy template to the output path:")))
			})
		})

		context("when templating fails", func() {
			it.Before(func() {
				templatizer.FillOutTemplateCall.Stub = nil
				templatizer.FillOutTemplateCall.Returns.Error = errors.New("some-error")
			})

			it("errors", func() {
				err := internal.Bootstrap(templatizer, "some-org/someBuildpack", outputPath, templatePath, executable)
				Expect(err).To(MatchError(`failed to fill out template: "some-error"`))
			})
		})

		context("when running go mod tidy fails", func() {
			it.Before(func() {
				executable.ExecuteCall.Returns.Error = errors.New("ooops, this blew up")
			})

			it("errors", func() {
				err := internal.Bootstrap(templatizer, "some-org/someBuildpack", outputPath, templatePath, executable)
				Expect(err).To(MatchError("failed to run 'go mod tidy': ooops, this blew up"))
			})
		})
	})
}
