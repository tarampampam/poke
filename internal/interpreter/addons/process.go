package addons

import (
	"os"
	"strings"

	js "github.com/dop251/goja"
)

type Process struct {
	Env map[string]string `json:"env"`
}

func NewProcess() *Process {
	var (
		environ = os.Environ()
		env     = make(map[string]string, len(environ))
	)

	for _, e := range environ {
		envKeyValue := strings.SplitN(e, "=", 2)

		if len(envKeyValue) == 2 {
			env[envKeyValue[0]] = envKeyValue[1]
		}
	}

	return &Process{Env: env}
}

func (p Process) Register(runtime *js.Runtime) error {
	return runtime.GlobalObject().DefineDataProperty(
		"process",
		runtime.ToValue(p),
		js.FLAG_FALSE, // writable
		js.FLAG_FALSE, // configurable
		js.FLAG_TRUE,  // enumerable
	)
}
