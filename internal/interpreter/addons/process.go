package addons

import (
	"os"
	"strings"

	js "github.com/dop251/goja"
)

type Process struct {
	env map[string]string
}

func NewProcess() *Process {
	return &Process{env: make(map[string]string)}
}

func (p Process) Register(runtime *js.Runtime) error {
	var process = runtime.NewObject()

	for _, e := range os.Environ() {
		envKeyValue := strings.SplitN(e, "=", 2)

		if len(envKeyValue) == 2 {
			p.env[envKeyValue[0]] = envKeyValue[1]
		}
	}

	if err := process.Set("env", p.env); err != nil {
		return err
	}

	return runtime.GlobalObject().DefineDataProperty(
		"process",
		process,
		js.FLAG_FALSE, // writable
		js.FLAG_FALSE, // configurable
		js.FLAG_TRUE,  // enumerable
	)
}
