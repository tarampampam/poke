package addons_test

import (
	"context"
	"errors"
	"testing"

	js "github.com/dop251/goja"
	"github.com/stretchr/testify/assert"

	"github.com/tarampampam/poke/internal/js/addons"
	"github.com/tarampampam/poke/internal/js/events"
)

func TestEvents_Push(t *testing.T) {
	var (
		runtime = js.New()
		channel = make(chan events.Event)
		addon   = addons.NewEvents(context.Background(), runtime, channel)
	)

	addon.Push(runtime.ToValue(map[string]any{}))

	event := <-channel

	assert.Equal(t, events.LevelDebug, event.Level)
	assert.Empty(t, event.Message)

	addon.Push(runtime.ToValue(map[string]any{"level": "info", "message": "foo1", "error": "bar"}))

	event = <-channel

	assert.Equal(t, events.LevelInfo, event.Level)
	assert.Equal(t, "foo1", event.Message)
	assert.Equal(t, errors.New("bar"), event.Error)
}

func TestEvents_Register(t *testing.T) {
	var (
		runtime = js.New()
		addon   = addons.NewEvents(context.Background(), runtime, nil)
	)

	const name = "events"

	assert.Nil(t, runtime.GlobalObject().Get(name))
	assert.NoError(t, addon.Register(runtime))
	assert.Same(t, addon, runtime.GlobalObject().Get(name).Export())
}
