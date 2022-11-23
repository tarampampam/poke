package addons

import (
	"encoding/base64"

	js "github.com/dop251/goja"
)

type Encoding struct {
	runtime *js.Runtime
}

func NewEncoding(runtime *js.Runtime) *Encoding { return &Encoding{runtime: runtime} }

func (e *Encoding) Base64encode(args ...js.Value) js.Value {
	if len(args) == 0 {
		return js.Undefined()
	}

	var encoder = base64.StdEncoding // std

	if len(args) > 1 {
		if modeOption := args[1].ToObject(e.runtime).Get("mode"); modeOption != nil && modeOption.String() == "url" {
			encoder = base64.URLEncoding // url
		}
	}

	return e.runtime.ToValue(
		encoder.EncodeToString([]byte(args[0].String())),
	)
}

func (e *Encoding) Base64decode(args ...js.Value) js.Value {
	if len(args) == 0 {
		return js.Undefined()
	}

	var encoder = base64.StdEncoding // std

	if len(args) > 1 {
		if modeOption := args[1].ToObject(e.runtime).Get("mode"); modeOption != nil && modeOption.String() == "url" {
			encoder = base64.URLEncoding // url
		}
	}

	s, err := encoder.DecodeString(args[0].String())
	if err != nil {
		return js.Undefined()
	}

	return e.runtime.ToValue(string(s))
}

func (e *Encoding) Register(runtime *js.Runtime) error {
	return runtime.GlobalObject().DefineDataProperty(
		"encoding",
		runtime.ToValue(e),
		js.FLAG_FALSE, // writable
		js.FLAG_FALSE, // configurable
		js.FLAG_TRUE,  // enumerable
	)
}
