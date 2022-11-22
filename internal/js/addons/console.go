package addons

import (
	"fmt"
	"reflect"
	"strconv"

	js "github.com/dop251/goja"
	"github.com/json-iterator/go"

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

func (c *Console) toJSON(in any) string {
	if in == nil {
		return "null"
	}

	j, err := c.json.Marshal(in)
	if err == nil {
		return string(j)
	}

	return fmt.Sprintf("cannot convert passed value to json (%s)", err.Error())
}

var (
	typePromise      = reflect.TypeOf((*js.Promise)(nil))
	typeSymbol       = reflect.TypeOf((*js.Symbol)(nil))
	typeFunctionCall = reflect.TypeOf((*js.FunctionCall)(nil))
)

func (c *Console) format(args []js.Value) (message string, extra []log.Extra) {
	if len(args) > 0 {
		if str, ok := args[0].Export().(string); ok {
			message = str
		}

		if len(args) > 1 {
			extra = make([]log.Extra, len(args)-1)

			for i, arg := range args[1:] {
				iStr := strconv.Itoa(i)

				switch arg.ExportType().Kind() {
				case typePromise.Kind():
					extra[i] = log.With(iStr, "<Promise>")

				case typeSymbol.Kind():
					extra[i] = log.With(iStr, "<Symbol>")

				case typeFunctionCall.Kind(), reflect.Func:
					extra[i] = log.With(iStr, "ƒ(…)")

				case reflect.Pointer, reflect.UnsafePointer:
					extra[i] = log.With(iStr, "<pointer>")

				default:
					extra[i] = log.With(iStr, c.toJSON(arg.Export()))
				}
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
