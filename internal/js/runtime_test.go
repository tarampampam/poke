package js_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tarampampam/poke/internal/js"
	"github.com/tarampampam/poke/internal/log"
)

func TestRuntime_RunScript(t *testing.T) {
	runtime, err := js.NewRuntime(context.Background(), log.NewNop())
	assert.NoError(t, err)

	defer runtime.Close()

	assert.NoError(t, runtime.RunScript("", "const a = 1;"))
}

func TestRuntime_EventPush(t *testing.T) {
	runtime, err := js.NewRuntime(context.Background(), log.NewNop())
	assert.NoError(t, err)

	defer runtime.Close()

	go func() { assert.NoError(t, runtime.RunScript("", "events.push({message: 'foo'})")) }()

	assert.Equal(t, "foo", (<-runtime.Events()).Message)
}

func TestRuntime_ConsoleMethods(t *testing.T) {
	runtime, err := js.NewRuntime(context.Background(), log.NewNop())
	assert.NoError(t, err)

	defer runtime.Close()

	assert.NoError(t, runtime.RunScript("", `console.debug('foo')
console.log('bar')
console.info('baz')
console.warn('qux')
console.error('quux')`))
}
