package interpreter

import (
	"context"
	_ "embed"

	js "github.com/dop251/goja"

	"github.com/tarampampam/poke/internal/interpreter/addons"
)

//go:embed global.js
var global string

type addonRegisterer interface {
	Register(*js.Runtime) error
}

type Runtime struct {
	js      *js.Runtime
	reports *addons.Reports
}

func NewRuntime(ctx context.Context) (*Runtime, error) {
	var runtime = Runtime{js: js.New()}

	runtime.reports = addons.NewReports(runtime.js)

	runtime.js.SetFieldNameMapper(js.TagFieldNameMapper("json", true))

	for _, addon := range []addonRegisterer{
		addons.NewIO(runtime.js),
		addons.NewProcess(),
		addons.NewFetch(ctx, nil),
		runtime.reports,
	} {
		if err := addon.Register(runtime.js); err != nil {
			return nil, err
		}
	}

	if _, err := runtime.js.RunScript("global", global); err != nil {
		return nil, err
	}

	return &runtime, nil
}

func (r *Runtime) RunString(script string) (any, error) { return r.js.RunString(script) }

func (r *Runtime) Reports() []addons.Report { return r.reports.Stack() }
