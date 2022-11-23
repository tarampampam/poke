package addons

import (
	"fmt"
	"reflect"
	"strconv"

	js "github.com/dop251/goja"
	jsoniter "github.com/json-iterator/go"

	"github.com/tarampampam/poke/internal/log"
)

type Console struct {
	runtime *js.Runtime
	log     log.Logger
	json    jsoniter.API
}

func NewConsole(runtime *js.Runtime, log log.Logger) *Console {
	return &Console{
		runtime: runtime,
		log:     log,
		json:    jsoniter.ConfigFastest,
	}
}

var (
	typePromise      = reflect.TypeOf((*js.Promise)(nil))      //nolint:gochecknoglobals
	typeFunctionCall = reflect.TypeOf((*js.FunctionCall)(nil)) //nolint:gochecknoglobals
)

func (c *Console) valueToString(v js.Value) string {
	if v == nil {
		return "null"
	} else if s, ok := v.Export().(string); ok {
		return s
	}

	switch {
	case js.IsUndefined(v):
		return "undefined"

	case js.IsNull(v):
		return "null"

	case js.IsNaN(v):
		return "NaN"

	case js.IsInfinity(v):
		return "Infinity"
	}

	if v.ExportType() == nil {
		return fmt.Sprintf("%v", v)
	}

	switch v.ExportType().Kind() { //nolint:exhaustive
	case typePromise.Kind():
		return "<Promise>"

	case typeFunctionCall.Kind(), reflect.Func:
		return "ƒ(…)"

	default:
		if j, err := c.json.Marshal(v); err == nil {
			return string(j)
		} else {
			return fmt.Sprintf("cannot convert passed value to json (%s)", err.Error())
		}
	}
}

func (c *Console) format(args []js.Value) (message string, extra []log.Extra) {
	if len(args) > 0 {
		if str, ok := args[0].Export().(string); ok {
			message = str
		} else if len(args) == 1 {
			message = c.valueToString(args[0])
		}

		if len(args) > 1 {
			extra = make([]log.Extra, len(args)-1)

			for i, arg := range args[1:] {
				extra[i] = log.With(strconv.Itoa(i), c.valueToString(arg))
			}
		}
	}

	return
}

func (c *Console) Debug(args ...js.Value) {
	if len(args) > 0 {
		msg, extra := c.format(args)

		c.log.Debug(msg, extra...)
	}
}

func (c *Console) Log(args ...js.Value) {
	if len(args) > 0 {
		msg, extra := c.format(args)

		c.log.Info(msg, extra...)
	}
}

func (c *Console) Info(args ...js.Value) {
	if len(args) > 0 {
		msg, extra := c.format(args)

		c.log.Info(msg, extra...)
	}
}

func (c *Console) Warn(args ...js.Value) {
	if len(args) > 0 {
		msg, extra := c.format(args)

		c.log.Warn(msg, extra...)
	}
}

func (c *Console) Error(args ...js.Value) {
	if len(args) > 0 {
		msg, extra := c.format(args)

		c.log.Error(msg, extra...)
	}
}

func (c *Console) Register(runtime *js.Runtime) error {
	return runtime.GlobalObject().DefineDataProperty(
		"console",
		runtime.ToValue(c),
		js.FLAG_FALSE, // writable
		js.FLAG_FALSE, // configurable
		js.FLAG_TRUE,  // enumerable
	)
}
