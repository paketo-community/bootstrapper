package bootstrapper_test

import (
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
		templatizer.FillOutTemplateCall.Stub = func(string, bootsrapper.Config) error {

	})

	it.After(func() {
		Expect(os.RemoveAll(outputPath)).To(Succeed())
		Expect(os.RemoveAll(configPath)).To(Succeed())
		Expect(os.RemoveAll(templatePath)).To(Succeed())
	})

	it.Focus("fills out every template in the correct path", func() {
		err := bootstrapper.Bootstrap(templatizer, configPath, templatePath, outputPath)
		Expect(err).NotTo(HaveOccurred())

		Expect(templatizer.FillOutTemplateCall.CallCount).To(Equal(2))
		Expect(templatizer.FillOutTemplateCall.Receives.Path).To(Equal(filepath.Join(outputPath, "templ2")))
		Expect(templatizer.FillOutTemplateCall.Receives.Config).To(Equal(bootstrapper.Config{
			Buildpack:    "someBuildpack",
			Organization: "some-org",
		}))
	})
}
