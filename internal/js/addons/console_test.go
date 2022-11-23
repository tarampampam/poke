package addons_test

import (
	"math"
	"testing"

	js "github.com/dop251/goja"
	"github.com/stretchr/testify/assert"

	"github.com/tarampampam/poke/internal/js/addons"
	"github.com/tarampampam/poke/internal/log"
)

type fakeLogger struct {
	LastMsg   string
	LastExtra []log.Extra
}

func (f *fakeLogger) Debug(msg string, v ...log.Extra)   { f.LastMsg, f.LastExtra = msg, v }
func (f *fakeLogger) Info(msg string, v ...log.Extra)    { f.LastMsg, f.LastExtra = msg, v }
func (f *fakeLogger) Success(msg string, v ...log.Extra) { f.LastMsg, f.LastExtra = msg, v }
func (f *fakeLogger) Warn(msg string, v ...log.Extra)    { f.LastMsg, f.LastExtra = msg, v }
func (f *fakeLogger) Error(msg string, v ...log.Extra)   { f.LastMsg, f.LastExtra = msg, v }
func (f *fakeLogger) Fatal(msg string, v ...log.Extra)   { f.LastMsg, f.LastExtra = msg, v }

var _ log.Logger = (*fakeLogger)(nil)

func TestConsole_TypesToString(t *testing.T) {
	var (
		runtime       = js.New()
		promise, _, _ = runtime.NewPromise()
	)

	for name, testCase := range map[string]struct {
		give js.Value
		want string
	}{
		"nan":       {js.NaN(), "NaN"},
		"false":     {runtime.ToValue(false), "false"},
		"true":      {runtime.ToValue(true), "true"},
		"null":      {js.Null(), "null"},
		"undefined": {js.Undefined(), "undefined"},
		"0":         {runtime.ToValue(0), "0"},
		"+infinity": {runtime.ToValue(math.Inf(+1)), "Infinity"},
		"-infinity": {runtime.ToValue(math.Inf(-1)), "Infinity"},
		"promise":   {runtime.ToValue(promise), "<Promise>"},
		"fn":        {runtime.ToValue(func(js.FunctionCall) js.Value { return js.NaN() }), "ƒ(…)"},
		"array":     {runtime.ToValue([]int{1, 2}), "[1,2]"},
		"object":    {runtime.ToValue(struct{ Foo, bar string }{"foo", "bar"}), `{"Foo":"foo"}`},
	} {
		tt := testCase

		t.Run(name, func(t *testing.T) {
			var (
				l     = &fakeLogger{}
				addon = addons.NewConsole(runtime, l)
			)

			addon.Log(tt.give)
			assert.Equal(t, l.LastMsg, tt.want)
		})
	}
}

func TestConsole_Methods(t *testing.T) {
	var (
		runtime, l = js.New(), &fakeLogger{}
		addon      = addons.NewConsole(runtime, l)
	)

	for name, testCase := range map[string]struct {
		giveHandler func(args ...js.Value)
		giveArgs    []js.Value
		wantMsg     string
		wantExtra   map[any]any
	}{
		"debug": {
			giveHandler: addon.Debug,
			giveArgs:    []js.Value{runtime.ToValue("foobar")},
			wantMsg:     "foobar",
		},
		"debug with args": {
			giveHandler: addon.Debug,
			giveArgs:    []js.Value{runtime.ToValue("foo"), runtime.ToValue("bar")},
			wantMsg:     "foo",
			wantExtra:   map[any]any{"0": "bar"},
		},
		"debug without args": {
			giveHandler: addon.Debug,
			giveArgs:    []js.Value{},
			wantMsg:     "",
		},
		"log": {
			giveHandler: addon.Log,
			giveArgs:    []js.Value{runtime.ToValue("foobar")},
			wantMsg:     "foobar",
		},
		"log with args": {
			giveHandler: addon.Log,
			giveArgs:    []js.Value{runtime.ToValue("foo"), runtime.ToValue("bar")},
			wantMsg:     "foo",
			wantExtra:   map[any]any{"0": "bar"},
		},
		"info": {
			giveHandler: addon.Info,
			giveArgs:    []js.Value{runtime.ToValue("foobar")},
			wantMsg:     "foobar",
		},
		"info with args": {
			giveHandler: addon.Info,
			giveArgs:    []js.Value{runtime.ToValue("foo"), runtime.ToValue("bar")},
			wantMsg:     "foo",
			wantExtra:   map[any]any{"0": "bar"},
		},
		"warn": {
			giveHandler: addon.Warn,
			giveArgs:    []js.Value{runtime.ToValue("foobar")},
			wantMsg:     "foobar",
		},
		"warn with args": {
			giveHandler: addon.Warn,
			giveArgs:    []js.Value{runtime.ToValue("foo"), runtime.ToValue("bar")},
			wantMsg:     "foo",
			wantExtra:   map[any]any{"0": "bar"},
		},
		"error": {
			giveHandler: addon.Error,
			giveArgs:    []js.Value{runtime.ToValue("foobar")},
			wantMsg:     "foobar",
		},
		"error with args": {
			giveHandler: addon.Error,
			giveArgs:    []js.Value{runtime.ToValue("foo"), runtime.ToValue("bar")},
			wantMsg:     "foo",
			wantExtra:   map[any]any{"0": "bar"},
		},
	} {
		tt := testCase

		t.Run(name, func(t *testing.T) {
			tt.giveHandler(tt.giveArgs...)

			var i = 0

			for key, value := range tt.wantExtra {
				assert.Equal(t, l.LastExtra[i].Key(), key)
				assert.Equal(t, l.LastExtra[i].Value(), value)

				i++
			}

			l.LastExtra, l.LastMsg = nil, "" // reset
		})
	}
}

func TestConsole_Register(t *testing.T) {
	var (
		runtime, l = js.New(), &fakeLogger{}
		addon      = addons.NewConsole(runtime, l)
	)

	const name = "console"

	assert.Nil(t, runtime.GlobalObject().Get(name))
	assert.NoError(t, addon.Register(runtime))
	assert.Same(t, addon, runtime.GlobalObject().Get(name).Export())
}
