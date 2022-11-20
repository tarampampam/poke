package js

import (
	"context"
	_ "embed"

	js "github.com/dop251/goja"

	"github.com/tarampampam/poke/internal/js/addons"
	"github.com/tarampampam/poke/internal/js/events"
)

//go:embed global.js
var global string

type addonRegisterer interface {
	Register(*js.Runtime) error
}

type Runtime struct {
	runtime *js.Runtime
	events  chan events.Event
}

func NewRuntime(ctx context.Context) (*Runtime, error) {
	var (
		r           = Runtime{runtime: js.New(), events: make(chan events.Event)}
		eventsAddon = addons.NewReports(r.runtime, r.events)
	)

	r.runtime.SetFieldNameMapper(js.TagFieldNameMapper("json", true))

	for _, addon := range []addonRegisterer{
		addons.NewIO(r.runtime),
		addons.NewProcess(),
		addons.NewFetch(ctx, nil),
		eventsAddon,
	} {
		if err := addon.Register(r.runtime); err != nil {
			return nil, err
		}
	}

	if _, err := r.runtime.RunScript("global.js", global); err != nil {
		return nil, err
	}

	return &r, nil
}

func (r *Runtime) Events() <-chan events.Event { return r.events }

func (r *Runtime) RunScript(name, script string) error {
	if _, err := r.runtime.RunScript(name, script); err != nil {
		return err
	}

	return nil
}

func (r *Runtime) Close() { close(r.events) }
