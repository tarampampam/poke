package addons_test

import (
	"context"
	"net/http"
	"testing"

	js "github.com/dop251/goja"
	"github.com/stretchr/testify/assert"

	"github.com/tarampampam/poke/internal/js/addons"
)

type httpClientFunc func(*http.Request) (*http.Response, error)

func (f httpClientFunc) Do(req *http.Request) (*http.Response, error) { return f(req) }

func TestFetch_Register(t *testing.T) {
	var (
		runtime = js.New()
		addon   = addons.NewFetch(context.Background(), nil)
	)

	const name = "fetchSync"

	assert.Nil(t, runtime.Get(name))
	assert.NoError(t, addon.Register(runtime))
	assert.NotNil(t, runtime.Get(name))
}
