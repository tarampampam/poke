package addons

import (
	"fmt"
	"os"
	"strings"

	js "github.com/dop251/goja"
)

// Console is an `console.log` alternative.
//
// The main idea was looked here: <https://github.com/dop251/goja/issues/396#issuecomment-1163556584>
type Console struct{}

func NewConsole() *Console { return &Console{} }

func (Console) formatForConsole(call ...js.Value) string {
	var output = make([]string, len(call))

	for i, arg := range call {
		output[i] = fmt.Sprintf("%v", arg.ToString())
	}

	return strings.Join(output, " ")
}

func (c Console) Log(args ...js.Value) {
	_, _ = fmt.Fprintln(os.Stdout, c.formatForConsole(args...))
}

func (c Console) Error(args ...js.Value) {
	_, _ = fmt.Fprintln(os.Stderr, c.formatForConsole(args...))
}

func (c Console) Debug(args ...js.Value) { c.Log(args...) }
func (c Console) Dir(args ...js.Value)   { c.Log(args...) }
func (c Console) Info(args ...js.Value)  { c.Log(args...) }
func (c Console) Warn(args ...js.Value)  { c.Log(args...) }
func (Console) Time(_ ...js.Value)       {} // doing nothing
func (Console) TimeEnd(_ ...js.Value)    {} // doing nothing
func (Console) Trace(_ ...js.Value)      {} // doing nothing
func (Console) Assert(_ ...js.Value)     {} // doing nothing

func (c Console) Register(runtime *js.Runtime) error {
	return runtime.GlobalObject().DefineDataProperty(
		"console",
		runtime.ToValue(c),
		js.FLAG_FALSE, // writable
		js.FLAG_FALSE, // configurable
		js.FLAG_TRUE,  // enumerable
	)
}
