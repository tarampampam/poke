package addons_test

import (
	"testing"

	js "github.com/dop251/goja"
	"github.com/stretchr/testify/assert"

	"github.com/tarampampam/poke/internal/js/addons"
)

func TestHashing_Md5(t *testing.T) {
	var (
		runtime = js.New()
		addon   = addons.NewHashing(runtime)
	)

	assert.True(t, js.IsUndefined(addon.Md5()))
	assert.Equal(t, "fc3ff98e8c6a0d3087d515c0473f8677", addon.Md5(runtime.ToValue("hello world!")).String())
	assert.Equal(t, "3858f62230ac3c915f300c664312c63f", addon.Md5(runtime.ToValue("foobar")).String())
}

func TestHashing_Sha256(t *testing.T) {
	var (
		runtime = js.New()
		addon   = addons.NewHashing(runtime)
	)

	assert.True(t, js.IsUndefined(addon.Md5()))
	assert.Equal(t,
		"7509e5bda0c762d2bac7f90d758b5b2263fa01ccbc542ab5e3df163be08e6ca9",
		addon.Sha256(runtime.ToValue("hello world!")).String(),
	)
	assert.Equal(t,
		"c3ab8ff13720e8ad9047dd39466b3c8974e592c2fa383d4a3960714caef0c4f2",
		addon.Sha256(runtime.ToValue("foobar")).String(),
	)
}

func TestHashing_Register(t *testing.T) {
	var (
		runtime = js.New()
		addon   = addons.NewHashing(runtime)
	)

	const name = "hashing"

	assert.Nil(t, runtime.GlobalObject().Get(name))
	assert.NoError(t, addon.Register(runtime))
	assert.Same(t, addon, runtime.GlobalObject().Get(name).Export())
}
