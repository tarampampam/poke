package addons

import (
	"io"

	js "github.com/dop251/goja"

	"github.com/tarampampam/poke/internal/js/printer"
)

type IO struct {
	runtime *js.Runtime
	stdOut  io.Writer
	stdErr  io.Writer
	printer printer.Printer
}

func NewIO(runtime *js.Runtime, stdOut, stdErr io.Writer, printer printer.Printer) *IO {
	return &IO{
		runtime: runtime,
		stdOut:  stdOut,
		stdErr:  stdErr,
		printer: printer,
	}
}

func (io *IO) write(w io.Writer, args ...js.Value) {
	var output = make([]any, len(args))

	for i, arg := range args {
		output[i] = arg
	}

	if err := io.printer(w, output...); err != nil {
		panic(io.runtime.ToValue(err.Error()))
	}
}

func (io *IO) StdOut(args ...js.Value) { io.write(io.stdOut, args...) }
func (io *IO) StdErr(args ...js.Value) { io.write(io.stdErr, args...) }

func (io *IO) Register(runtime *js.Runtime) error {
	return runtime.GlobalObject().DefineDataProperty(
		"io",
		runtime.ToValue(io),
		js.FLAG_FALSE, // writable
		js.FLAG_FALSE, // configurable
		js.FLAG_TRUE,  // enumerable
	)
}
