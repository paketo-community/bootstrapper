package bootstrapper

import (
	"fmt"
	"html/template"
	"os"
	"strings"
)

type Templatizer struct{}

func NewTemplatizer() Templatizer {
	return Templatizer{}
}

func RemoveHyphens(s string) string {
	return strings.ReplaceAll(s, "-", "")
}

func (tz Templatizer) FillOutTemplate(path string, config Config) error {
	templ, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	// go.mod file: update the module to use templating so the go.mod will be functional
	if strings.Contains(path, "go.mod") {
		newContents := strings.Replace(string(templ), "github.com/test/test", "github.com/{{ .Organization }}/{{ .Buildpack }}", -1)
		templ = []byte(newContents)
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("failed to open template file: %w", err)
	}

	funcMap := template.FuncMap{
		"Title":         strings.Title,
		"RemoveHyphens": RemoveHyphens,
	}
	t := template.Must(template.New("t1").Funcs(funcMap).Parse(string(templ)))

	err = t.Execute(file, config)
	if err != nil {
		return fmt.Errorf("failed to fill out template: %w", err)
	}

	file.Close()

	return nil
}
