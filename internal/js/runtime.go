package js

import (
	"context"
	_ "embed"
	"sync"

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

	closeOnce sync.Once
}

func NewRuntime(ctx context.Context) (*Runtime, error) {
	var r = Runtime{runtime: js.New(), events: make(chan events.Event, 32)} //nolint:gomnd

	r.runtime.SetFieldNameMapper(js.TagFieldNameMapper("json", true))

	for _, addon := range []addonRegisterer{
		addons.NewIO(r.runtime),
		addons.NewProcess(),
		addons.NewFetch(ctx, nil),
		addons.NewEvents(ctx, r.runtime, r.events),
	} {
		if err := addon.Register(r.runtime); err != nil {
			r.Close()

			return nil, err
		}
	}

	if _, err := r.runtime.RunScript("global.js", global); err != nil {
		r.Close()

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

func (r *Runtime) Interrupt(reason string) {
	r.runtime.Interrupt(reason)
}

func (r *Runtime) Close() {
	r.closeOnce.Do(func() {
		close(r.events)
	})
}
