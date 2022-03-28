package fakes

import (
	"sync"

	"github.com/paketo-community/bootstrapper/commands/internal"
)

type TemplateWriter struct {
	FillOutTemplateCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Path   string
			Config internal.Config
		}
		Returns struct {
			Error error
		}
		Stub func(string, internal.Config) error
	}
}

func (f *TemplateWriter) FillOutTemplate(param1 string, param2 internal.Config) error {
	f.FillOutTemplateCall.mutex.Lock()
	defer f.FillOutTemplateCall.mutex.Unlock()
	f.FillOutTemplateCall.CallCount++
	f.FillOutTemplateCall.Receives.Path = param1
	f.FillOutTemplateCall.Receives.Config = param2
	if f.FillOutTemplateCall.Stub != nil {
		return f.FillOutTemplateCall.Stub(param1, param2)
	}
	return f.FillOutTemplateCall.Returns.Error
}
