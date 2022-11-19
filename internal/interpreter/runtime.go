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
	Register(runtime *js.Runtime) error
}

func NewRuntime() (*js.Runtime, error) {
	var runtime = js.New()

	runtime.SetFieldNameMapper(js.TagFieldNameMapper("json", true))

	for _, addon := range []addonRegisterer{
		addons.NewProcess(),
		addons.NewConsole(),
		addons.NewFetch(context.TODO(), nil),
	} {
		if err := addon.Register(runtime); err != nil {
			return nil, err
		}
	}

	if _, err := runtime.RunScript("global", global); err != nil {
		return nil, err
	}

	return runtime, nil
}
