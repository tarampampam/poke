package addons

import (
	"context"
	"errors"
	"strings"

	js "github.com/dop251/goja"

	"github.com/tarampampam/poke/internal/js/events"
)

type Events struct {
	ctx     context.Context     `js:"-"`
	runtime *js.Runtime         `js:"-"`
	channel chan<- events.Event `js:"-"`
}

func NewEvents(ctx context.Context, runtime *js.Runtime, channel chan<- events.Event) *Events {
	return &Events{
		ctx:     ctx,
		runtime: runtime,
		channel: channel,
	}
}

func (e *Events) Push(args ...js.Value) {
	for _, arg := range args {
		var (
			obj   = arg.ToObject(e.runtime)
			event = events.Event{Level: events.LevelDebug}
		)

		if lvl := obj.Get("level"); lvl != nil {
			if value, isString := lvl.Export().(string); isString {
				event.Level = events.Level(strings.ToLower(value))
			}
		}

		if msg := obj.Get("message"); msg != nil {
			if value, isString := msg.Export().(string); isString {
				event.Message = value
			}
		}

		if err := obj.Get("error"); err != nil {
			event.Error = errors.New(err.String())
		}

		select {
		case <-e.ctx.Done():
			return

		case e.channel <- event:
		}
	}
}

func (e *Events) Register(runtime *js.Runtime) error {
	return runtime.GlobalObject().DefineDataProperty(
		"events",
		runtime.ToValue(e),
		js.FLAG_FALSE, // writable
		js.FLAG_FALSE, // configurable
		js.FLAG_TRUE,  // enumerable
	)
}
