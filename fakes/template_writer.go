package fakes

import (
	"sync"

	"github.com/paketo-community/bootstrapper"
)

type TemplateWriter struct {
	FillOutTemplateCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			Path   string
			Config bootstrapper.Config
		}
		Returns struct {
			Error error
		}
		Stub func(string, bootstrapper.Config) error
	}
}

func (f *TemplateWriter) FillOutTemplate(param1 string, param2 bootstrapper.Config) error {
	f.FillOutTemplateCall.Lock()
	defer f.FillOutTemplateCall.Unlock()
	f.FillOutTemplateCall.CallCount++
	f.FillOutTemplateCall.Receives.Path = param1
	f.FillOutTemplateCall.Receives.Config = param2
	if f.FillOutTemplateCall.Stub != nil {
		return f.FillOutTemplateCall.Stub(param1, param2)
	}
	return f.FillOutTemplateCall.Returns.Error
}
