package addons

import (
	"fmt"
	"io"
	"os"

	js "github.com/dop251/goja"
)

type IO struct {
	runtime *js.Runtime `js:"-"` // runtime instance
}

func NewIO(runtime *js.Runtime) *IO { return &IO{runtime: runtime} }

func (io IO) write(w io.Writer, args ...js.Value) {
	var output = make([]any, len(args))

	for i, arg := range args {
		output[i] = arg
	}

	if _, err := fmt.Fprint(w, output...); err != nil {
		panic(io.runtime.ToValue(err.Error()))
	}
}

func (io IO) StdOut(args ...js.Value) { io.write(os.Stdout, args...) }
func (io IO) StdErr(args ...js.Value) { io.write(os.Stderr, args...) }

func (io IO) Register(runtime *js.Runtime) error {
	return runtime.GlobalObject().DefineDataProperty(
		"io",
		runtime.ToValue(io),
		js.FLAG_FALSE, // writable
		js.FLAG_FALSE, // configurable
		js.FLAG_TRUE,  // enumerable
	)
}
