package interpreter

import (
	_ "embed"

	js "github.com/dop251/goja"

	"github.com/tarampampam/poke/internal/interpreter/addons"
)

//go:embed helpers.js
var helpers string

type addonRegisterer interface {
	Register(runtime *js.Runtime) error
}

func NewRuntime() (*js.Runtime, error) {
	var runtime = js.New()

	for _, addon := range []addonRegisterer{
		addons.NewProcess(),
		&addons.Console{},
	} {
		if err := addon.Register(runtime); err != nil {
			return nil, err
		}
	}

	if _, err := runtime.RunScript("helpers", helpers); err != nil {
		return nil, err
	}

	return runtime, nil
}
