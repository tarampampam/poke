package addons

import (
	"io"
	"os"

	js "github.com/dop251/goja"

	"github.com/tarampampam/poke/internal/js/printer"
)

type IO struct {
	runtime  *js.Runtime
	printer  printer.Printer
	logLevel string
}

func NewIO(runtime *js.Runtime, printer printer.Printer, logLevel string) *IO {
	var lvl = "info"

	switch logLevel {
	case "debug", "verbose":
		lvl = "debug"

	case "info":
		lvl = "info"

	case "warn", "warning":
		lvl = "warning"

	case "err", "error", "fatal":
		lvl = "error"
	}

	return &IO{runtime: runtime, printer: printer, logLevel: lvl}
}

func (io IO) write(w io.Writer, args ...js.Value) {
	var output = make([]any, len(args))

	for i, arg := range args {
		output[i] = arg
	}

	if err := io.printer(w, output...); err != nil {
		panic(io.runtime.ToValue(err.Error()))
	}
}

func (io IO) StdOut(args ...js.Value) { io.write(os.Stdout, args...) }
func (io IO) StdErr(args ...js.Value) { io.write(os.Stderr, args...) }
func (io IO) LogLevel() string        { return io.logLevel }

func (io IO) Register(runtime *js.Runtime) error {
	return runtime.GlobalObject().DefineDataProperty(
		"io",
		runtime.ToValue(io),
		js.FLAG_FALSE, // writable
		js.FLAG_FALSE, // configurable
		js.FLAG_TRUE,  // enumerable
	)
}
