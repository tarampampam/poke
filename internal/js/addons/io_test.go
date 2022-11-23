package addons_test

import (
	"bytes"
	"testing"

	js "github.com/dop251/goja"
	"github.com/stretchr/testify/assert"

	"github.com/tarampampam/poke/internal/js/addons"
	"github.com/tarampampam/poke/internal/js/printer"
)

func TestIO_StdOut(t *testing.T) {
	var (
		runtime        = js.New()
		stdOut, stdErr bytes.Buffer
		addon          = addons.NewIO(runtime, &stdOut, &stdErr, printer.DefaultPrinter())
	)

	addon.StdOut(runtime.ToValue("foo"))

	assert.Equal(t, "foo", stdOut.String())
	assert.Empty(t, stdErr.String())

	stdOut.Reset()
	stdErr.Reset()

	addon.StdErr(runtime.ToValue("bar"))

	assert.Empty(t, stdOut.String())
	assert.Equal(t, "bar", stdErr.String())
}

func TestIO_Register(t *testing.T) {
	var (
		runtime = js.New()
		buf     bytes.Buffer
		addon   = addons.NewIO(runtime, &buf, &buf, printer.DefaultPrinter())
	)

	const name = "io"

	assert.Nil(t, runtime.GlobalObject().Get(name))
	assert.NoError(t, addon.Register(runtime))
	assert.Equal(t, addon, runtime.GlobalObject().Get(name).Export())
}
