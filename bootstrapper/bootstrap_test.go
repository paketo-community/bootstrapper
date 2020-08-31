package bootstrapper_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-community/bootstrapper/bootstrapper"
	"github.com/paketo-community/bootstrapper/bootstrapper/fakes"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testBootstrap(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		templatizer *fakes.TemplateWriter

		outputPath   string
		templatePath string
		err          error

		templatizerCalls []struct {
			path   string
			config bootstrapper.Config
		}
	)

	it.Before(func() {
		outputPath, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		templatePath, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		err = ioutil.WriteFile(filepath.Join(templatePath, "templ1"), []byte(""), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		err = ioutil.WriteFile(filepath.Join(templatePath, "templ2"), []byte(""), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		templatizer = &fakes.TemplateWriter{}
		templatizer.FillOutTemplateCall.Stub = func(path string, config bootstrapper.Config) error {
			templatizerCalls = append(templatizerCalls, struct {
				path   string
				config bootstrapper.Config
			}{path: path, config: config})
			return nil
		}
	})

	it.After(func() {
		Expect(os.RemoveAll(outputPath)).To(Succeed())
		Expect(os.RemoveAll(templatePath)).To(Succeed())
	})

	it("fills out every template in the correct path", func() {
		err := bootstrapper.Bootstrap(templatizer, "some-org/someBuildpack", templatePath, outputPath)
		Expect(err).NotTo(HaveOccurred())

		Expect(templatizer.FillOutTemplateCall.CallCount).To(Equal(2))
		Expect(templatizerCalls[0].path).To(Equal(filepath.Join(outputPath, "templ1")))
		Expect(templatizerCalls[0].config).To(Equal(bootstrapper.Config{
			Buildpack:    "someBuildpack",
			Organization: "some-org",
		}))

		Expect(templatizerCalls[1].path).To(Equal(filepath.Join(outputPath, "templ2")))
		Expect(templatizerCalls[1].config).To(Equal(bootstrapper.Config{
			Buildpack:    "someBuildpack",
			Organization: "some-org",
		}))
	})

	context("error cases", func() {
		context("when the buildpack name is malformed", func() {
			it("errors with a helpful message", func() {
				err := bootstrapper.Bootstrap(templatizer, "some-malformed-name", templatePath, outputPath)
				Expect(err).To(HaveOccurred())

				Expect(err).To(MatchError(ContainSubstring("buildpack name must be in format <organization>/<buildpack-name>")))
			})
		})
		context("when the template can not be copied to the output path", func() {
			it.Before(func() {
				templatePath = "does-not-exist"
			})

			it("errors", func() {
				err := bootstrapper.Bootstrap(templatizer, "some-org/someBuildpack", templatePath, outputPath)
				Expect(err).To(HaveOccurred())

				Expect(err).To(MatchError(ContainSubstring("failed to copy template to the output path:")))
			})
		})

		context("when templating fails", func() {
			it.Before(func() {
				templatizer.FillOutTemplateCall.Stub = nil
				templatizer.FillOutTemplateCall.Returns.Error = errors.New("some-error")
			})

			it("errors", func() {
				err := bootstrapper.Bootstrap(templatizer, "some-org/someBuildpack", templatePath, outputPath)
				Expect(err).To(HaveOccurred())

				Expect(err).To(MatchError(`failed to fill out template: "some-error"`))
			})
		})
	})
}
