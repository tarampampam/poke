package addons

import (
	"context"
	"os"
	"strings"
	"time"

	js "github.com/dop251/goja"
)

type Process struct {
	ctx     context.Context
	runtime *js.Runtime

	Env map[string]string `json:"env"`
}

func NewProcess(ctx context.Context, runtime *js.Runtime) *Process {
	var (
		environ = os.Environ()
		env     = make(map[string]string, len(environ))
	)

	for _, e := range environ {
		envKeyValue := strings.SplitN(e, "=", 2) //nolint:gomnd

		if len(envKeyValue) == 2 { //nolint:gomnd
			env[envKeyValue[0]] = envKeyValue[1]
		}
	}

	return &Process{ctx: ctx, runtime: runtime, Env: env}
}

// Delay is a helper function for delaying script execution for a given duration.
func (p Process) Delay(args ...js.Value) {
	if len(args) == 0 || p.ctx.Err() != nil {
		return
	}

	if delay := args[0].ToInteger(); delay > 0 {
		var t = time.NewTimer(time.Duration(delay) * time.Millisecond)

		defer t.Stop()

		select {
		case <-p.ctx.Done():
			return

		case <-t.C: // do nothing
		}
	}
}

// Interrupt is a helper function for interrupting script execution.
func (p Process) Interrupt(args ...js.Value) {
	var reason = "interrupted by the script"

	if len(args) == 1 {
		reason = args[0].String()
	}

	p.runtime.Interrupt(reason)
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
