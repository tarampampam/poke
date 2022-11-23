package addons_test

import (
	"testing"

	js "github.com/dop251/goja"
	"github.com/stretchr/testify/assert"

	"github.com/tarampampam/poke/internal/js/addons"
)

func TestFaker_Methods(t *testing.T) { //nolint:gocyclo
	var (
		runtime = js.New()
		addon   = addons.NewFaker(runtime)
	)

	assert.NotPanics(t, func() { addon.Bool() })

	for i := 0; i < 100; i++ {
		var falsy = addon.Falsy()

		assert.True(t,
			js.IsUndefined(falsy) || js.IsNaN(falsy) || js.IsNull(falsy) ||
				falsy.Export() == false || falsy.Export() == "" || falsy.Export() == int64(0),
		)
	}

	{
		var character = addon.Character(runtime.ToValue(map[string]any{"pool": "abc"}))

		assert.True(t, character == "a" || character == "b" || character == "c")
		assert.NotEmpty(t, addon.Character())
	}

	for i := 0; i < 100; i++ {
		var floating = addon.Floating()

		assert.GreaterOrEqual(t, floating, float32(-32768))
		assert.LessOrEqual(t, floating, float32(32768))
	}

	for i := 0; i < 100; i++ {
		var integer = addon.Integer()

		assert.GreaterOrEqual(t, integer, -9007199254740991)
		assert.LessOrEqual(t, integer, 9007199254740991)

		integer = addon.Integer(runtime.ToValue(map[string]any{"min": -100, "max": -90}))

		assert.GreaterOrEqual(t, integer, -100)
		assert.LessOrEqual(t, integer, -90)
	}

	for i := 0; i < 100; i++ {
		var letter = addon.Letter()

		assert.Len(t, letter, 1)
		assert.Contains(t, "abcdefghijklmnopqrstuvwxyz", letter)
	}

	for i := 0; i < 100; i++ {
		var str = addon.String()

		assert.Len(t, str, 11)

		for _, r := range str {
			assert.Contains(t, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", string(r))
		}

		str = addon.String(runtime.ToValue(map[string]any{"pool": "abc", "length": 3}))

		assert.Len(t, str, 3)

		for _, r := range str {
			assert.Contains(t, "abc", string(r))
		}
	}

	assert.NotEmpty(t, addon.Paragraph())
	assert.NotEmpty(t, addon.Word())
	assert.NotEmpty(t, addon.Domain())
	assert.NotEmpty(t, addon.Email())
	assert.NotEmpty(t, addon.Ip())
	assert.NotEmpty(t, addon.Ipv6())
	assert.NotEmpty(t, addon.Tld())
	assert.NotEmpty(t, addon.Url())
	assert.NotNil(t, addon.Date())

	for i := 0; i < 100; i++ {
		const hashSet = "0123456789abcdef"

		var hash = addon.Hash()

		assert.Len(t, hash, 40)

		for _, r := range hash {
			assert.Contains(t, hashSet, string(r))
		}

		hash = addon.Hash(runtime.ToValue(map[string]any{"length": 3}))

		assert.Len(t, hash, 3)

		for _, r := range hash {
			assert.Contains(t, hashSet, string(r))
		}
	}

	assert.NotEmpty(t, addon.Uuid())

	{
		r1, r2, r3 := runtime.ToValue(1), runtime.ToValue(2), runtime.ToValue(3)

		var pick = addon.Random(r1, r2, r3)

		assert.True(t, pick == r1 || pick == r2 || pick == r3)

		assert.True(t, js.IsUndefined(addon.Random()))
	}
}

func TestFaker_Register(t *testing.T) {
	var (
		runtime = js.New()
		addon   = addons.NewFaker(runtime)
	)

	const name = "faker"

	assert.Nil(t, runtime.GlobalObject().Get(name))
	assert.NoError(t, addon.Register(runtime))
	assert.Equal(t, addon, runtime.GlobalObject().Get(name).Export())
}
