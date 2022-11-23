package addons_test

import (
	"bytes"
	"math"
	"testing"

	js "github.com/dop251/goja"
	"github.com/stretchr/testify/assert"

	"github.com/tarampampam/poke/internal/js/addons"
	"github.com/tarampampam/poke/internal/log"
)

func TestConsole_TypesToString(t *testing.T) {
	var runtime = js.New()

	var promise, _, _ = runtime.NewPromise()

	for name, tt := range map[string]struct {
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
		tt := tt
		t.Run(name, func(t *testing.T) {
			var (
				buf   bytes.Buffer
				addon = addons.NewConsole(runtime, log.New(log.DebugLevel, log.WithStdOut(&buf), log.WithoutPrefix()))
			)

			addon.Log(tt.give)
			assert.Equal(t, buf.String(), tt.want+"\n")

			buf.Reset()
		})
	}
}
