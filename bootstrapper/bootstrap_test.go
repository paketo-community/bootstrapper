package bootstrapper_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/joshzarrabi/cnb-bootstrapper/bootstrapper"
	"github.com/joshzarrabi/cnb-bootstrapper/bootstrapper/fakes"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testBootstrap(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		templatizer *fakes.TemplateWriter

		outputPath   string
		configPath   string
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

		configFile, err := ioutil.TempFile("", "config.yml")
		Expect(err).NotTo(HaveOccurred())

		_, err = configFile.WriteString(`---
organization: some-org
buildpack: someBuildpack
`)
		Expect(err).NotTo(HaveOccurred())

		templatePath, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		err = ioutil.WriteFile(filepath.Join(templatePath, "templ1"), []byte(""), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		err = ioutil.WriteFile(filepath.Join(templatePath, "templ2"), []byte(""), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		configPath = configFile.Name()

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
		Expect(os.RemoveAll(configPath)).To(Succeed())
		Expect(os.RemoveAll(templatePath)).To(Succeed())
	})

	it("fills out every template in the correct path", func() {
		err := bootstrapper.Bootstrap(templatizer, configPath, templatePath, outputPath)
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
		context("when the config file can not be read", func() {
			it.Before(func() {
				Expect(os.RemoveAll(configPath)).To(Succeed())
			})

			it("errors", func() {
				err := bootstrapper.Bootstrap(templatizer, configPath, templatePath, outputPath)
				Expect(err).To(HaveOccurred())

				Expect(err).To(MatchError(ContainSubstring("failed to read config file:")))
			})
		})

		context("when the config file can not be parsed", func() {
			it.Before(func() {
				configFile, err := ioutil.TempFile("", "config.yml")
				Expect(err).NotTo(HaveOccurred())

				_, err = configFile.WriteString(`some-bad-yaml`)
				Expect(err).NotTo(HaveOccurred())

				configPath = configFile.Name()
			})

			it("errors", func() {
				err := bootstrapper.Bootstrap(templatizer, configPath, templatePath, outputPath)
				Expect(err).To(HaveOccurred())

				Expect(err).To(MatchError(ContainSubstring("failed to parse config file:")))
			})
		})

		context("when the template can not be copied to the output path", func() {
			it.Before(func() {
				templatePath = "does-not-exist"
			})

			it("errors", func() {
				err := bootstrapper.Bootstrap(templatizer, configPath, templatePath, outputPath)
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
				err := bootstrapper.Bootstrap(templatizer, configPath, templatePath, outputPath)
				Expect(err).To(HaveOccurred())

				Expect(err).To(MatchError(`failed to fill out template: "some-error"`))
			})
		})
	})
}
