package addons_test

import (
	"testing"

	js "github.com/dop251/goja"
	"github.com/stretchr/testify/assert"

	"github.com/tarampampam/poke/internal/js/addons"
)

func TestEncoding_Base64EncodeAndDecode(t *testing.T) {
	var (
		runtime = js.New()
		addon   = addons.NewEncoding(runtime)
	)

	assert.True(t, js.IsUndefined(addon.Base64decode()))
	assert.True(t, js.IsUndefined(addon.Base64encode()))

	const decoded, encoded = "Hello world!", "SGVsbG8gd29ybGQh"

	assert.Equal(t, encoded, addon.Base64encode(runtime.ToValue(decoded)).String())
	assert.Equal(t, decoded, addon.Base64decode(runtime.ToValue(encoded)).String())

	assert.True(t, js.IsUndefined(addon.Base64decode(runtime.ToValue("foobar"))))
}

func TestEncoding_Register(t *testing.T) {
	var (
		runtime = js.New()
		addon   = addons.NewEncoding(runtime)
	)

	const name = "encoding"

	assert.Nil(t, runtime.GlobalObject().Get(name))
	assert.NoError(t, addon.Register(runtime))
	assert.Same(t, addon, runtime.GlobalObject().Get(name).Export())
}
