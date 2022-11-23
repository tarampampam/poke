package addons_test

import (
	"context"
	"os"
	"testing"
	"time"

	js "github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tarampampam/poke/internal/js/addons"
)

func TestProcessEnv(t *testing.T) {
	require.NoError(t, os.Setenv("TEST_ENV", "test"))

	defer func() { _ = os.Unsetenv("TEST_ENV") }()

	var (
		runtime = js.New()
		addon   = addons.NewProcess(context.Background(), runtime)
	)

	assert.NotEmpty(t, addon.Env)
	assert.Equal(t, "test", addon.Env["TEST_ENV"])
}

func TestProcess_Delay(t *testing.T) {
	var (
		runtime        = js.New()
		addon          = addons.NewProcess(context.Background(), runtime)
		startedAtMilli = time.Now().UnixMilli()
	)

	addon.Delay(runtime.ToValue(50))

	var endedAtMilli = time.Now().UnixMilli()

	assert.GreaterOrEqual(t, startedAtMilli+53, endedAtMilli)
	assert.LessOrEqual(t, endedAtMilli, startedAtMilli+100) // 50 millis is allowed delta
}

func TestProcess_Interrupt(t *testing.T) {
	var (
		runtime = js.New()
		addon   = addons.NewProcess(context.Background(), runtime)
	)

	addon.Interrupt(runtime.ToValue("foo"))
	addon.Interrupt(runtime.ToValue("bar")) // multiple

	_, err := runtime.RunString("")
	assert.ErrorContains(t, err, "bar")
}

func TestProcess_Register(t *testing.T) {
	var (
		runtime = js.New()
		addon   = addons.NewProcess(context.Background(), runtime)
	)

	const name = "process"

	assert.Nil(t, runtime.GlobalObject().Get(name))
	assert.NoError(t, addon.Register(runtime))
	assert.Equal(t, addon, runtime.GlobalObject().Get(name).Export())
}
