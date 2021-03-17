package bootstrapper_test

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/paketo-community/bootstrapper/bootstrapper"
	"github.com/sclevine/spec"
)

func testTemplatizer(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		config       bootstrapper.Config
		templatePath string
		templatizer  bootstrapper.Templatizer
	)

	context("FillOutTemplate", func() {
		it.Before(func() {
			config = bootstrapper.Config{
				Buildpack:    "my-bp-name",
				Organization: "myorg",
			}

			template, err := ioutil.TempFile("", "template")
			Expect(err).NotTo(HaveOccurred())

			_, err = template.WriteString("{{ .Buildpack }} {{ .Organization }}")
			Expect(err).NotTo(HaveOccurred())

			templatePath = template.Name()

			templatizer = bootstrapper.NewTemplatizer()
		})

		it.After(func() {
			Expect(os.RemoveAll(templatePath)).To(Succeed())
		})

		it("fills out the template based on the config", func() {
			err := templatizer.FillOutTemplate(templatePath, config)
			Expect(err).NotTo(HaveOccurred())

			contents, err := ioutil.ReadFile(templatePath)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(contents)).To(Equal("my-bp-name myorg"))
		})

		context("when the templated file is a go.mod", func() {
			it.Before(func() {
				Expect(os.RemoveAll(templatePath)).To(Succeed())

				template, err := ioutil.TempFile("", "go.mod")
				Expect(err).NotTo(HaveOccurred())

				_, err = template.WriteString("module github.com/test/test")
				Expect(err).NotTo(HaveOccurred())

				templatePath = template.Name()
			})

			it("templates the module", func() {
				err := templatizer.FillOutTemplate(templatePath, config)
				Expect(err).NotTo(HaveOccurred())

				contents, err := ioutil.ReadFile(templatePath)
				Expect(err).NotTo(HaveOccurred())

				Expect(string(contents)).To(Equal("module github.com/myorg/my-bp-name"))
			})
		})

		context("when the template file uses the RemoveHyphens function", func() {
			it.Before(func() {
				Expect(os.RemoveAll(templatePath)).To(Succeed())

				template, err := ioutil.TempFile("", "template")
				Expect(err).NotTo(HaveOccurred())

				_, err = template.WriteString("{{ .Buildpack | RemoveHyphens }}")
				Expect(err).NotTo(HaveOccurred())

				templatePath = template.Name()
			})

			it("removes hyphens from the templated value", func() {
				err := templatizer.FillOutTemplate(templatePath, config)
				Expect(err).NotTo(HaveOccurred())

				contents, err := ioutil.ReadFile(templatePath)
				Expect(err).NotTo(HaveOccurred())

				Expect(string(contents)).To(Equal("mybpname"))
			})
		})

		context("error cases", func() {
			context("when the template file can not be read", func() {
				it.Before(func() {
					Expect(os.RemoveAll(templatePath)).To(Succeed())
				})

				it("errors", func() {
					err := templatizer.FillOutTemplate(templatePath, config)
					Expect(err).To(HaveOccurred())

					Expect(err.Error()).To(ContainSubstring("failed to read template file"))
				})
			})

			context("when the file can not be opened", func() {
				it.Before(func() {
					Expect(os.Chmod(templatePath, 0444)).To(Succeed())
				})

				it("errors", func() {
					err := templatizer.FillOutTemplate(templatePath, config)
					Expect(err).To(HaveOccurred())

					Expect(err.Error()).To(ContainSubstring("failed to open template file"))
				})
			})

			context("when the template can not be filled out", func() {
				it.Before(func() {
					template, err := ioutil.TempFile("", "template")
					Expect(err).NotTo(HaveOccurred())

					_, err = template.WriteString("{{ .buildpack }} ")
					Expect(err).NotTo(HaveOccurred())

					templatePath = template.Name()
				})

				it("errors", func() {
					err := templatizer.FillOutTemplate(templatePath, config)
					Expect(err).To(HaveOccurred())

					Expect(err.Error()).To(ContainSubstring("failed to fill out template"))
				})

			})
		})
	})
}
