package js

import (
	"context"
	_ "embed"
	"os"
	"sync"

	js "github.com/dop251/goja"
	"github.com/dop251/goja/parser"
	"github.com/pkg/errors"

	"github.com/tarampampam/poke/internal/js/addons"
	"github.com/tarampampam/poke/internal/js/events"
	"github.com/tarampampam/poke/internal/js/printer"
	"github.com/tarampampam/poke/internal/log"
)

//go:embed global.js
var global string

//go:embed global.d.ts
var globalDts string

func DTS() string { return globalDts }

type (
	// RuntimeOption allows to set up some internal Runtime properties from outside.
	RuntimeOption func(*Runtime)

	// Runtime is a wrapper for goja.Runtime.
	Runtime struct {
		runtime *js.Runtime
		events  chan events.Event
		printer printer.Printer

		closeOnce sync.Once
	}

	// addonRegisterer is an interface for all addons.
	addonRegisterer interface {
		Register(*js.Runtime) error
	}
)

// WithPrinter sets up the printer for the runtime.
func WithPrinter(p printer.Printer) RuntimeOption {
	return func(r *Runtime) { r.printer = p }
}

// NewRuntime creates new Runtime instance. Don't forget to close it after usage.
func NewRuntime(ctx context.Context, log log.Logger, options ...RuntimeOption) (*Runtime, error) {
	var r = &Runtime{ // defaults
		runtime: js.New(),
		events:  make(chan events.Event, 32), //nolint:gomnd
		printer: printer.DefaultPrinter(),
	}

	r.runtime.SetFieldNameMapper(js.TagFieldNameMapper("json", true))
	r.runtime.SetParserOptions(parser.WithDisableSourceMaps)

	for _, opt := range options {
		opt(r)
	}

	for _, addon := range []addonRegisterer{
		addons.NewIO(r.runtime, os.Stdout, os.Stderr, r.printer),
		addons.NewConsole(r.runtime, log),
		addons.NewProcess(ctx, r.runtime),
		addons.NewFetch(ctx, nil),
		addons.NewEvents(ctx, r.runtime, r.events),
		addons.NewFaker(r.runtime),
		addons.NewEncoding(r.runtime),
		addons.NewHashing(r.runtime),
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

	return r, nil
}

// Events returns channel with events. Channel reading is required for the events working.
func (r *Runtime) Events() <-chan events.Event { return r.events }

// RunScript runs the JS script.
func (r *Runtime) RunScript(name, script string) error {
	if _, err := r.runtime.RunScript(name, script); err != nil {
		return err
	}

	if afterScript, ok := js.AssertFunction(r.runtime.Get("init")); ok {
		if _, err := afterScript(r.runtime.GlobalObject()); err != nil {
			return errors.Wrap(err, "init() calling failed")
		}
	}

	return nil
}

// Interrupt interrupts the runtime.
func (r *Runtime) Interrupt(reason string) {
	r.runtime.Interrupt(reason)
}

// Close closes the runtime.
func (r *Runtime) Close() {
	r.closeOnce.Do(func() {
		close(r.events)
	})
}
