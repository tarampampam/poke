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
	runtime, _ := js.NewRuntime(context.Background(), log.NewNop())

	defer runtime.Close()

	go func() { assert.NoError(t, runtime.RunScript("", "events.push({message: 'foo'})")) }()

	assert.Equal(t, "foo", (<-runtime.Events()).Message)
}

func TestRuntime_ConsoleMethods(t *testing.T) {
	runtime, _ := js.NewRuntime(context.Background(), log.NewNop())

	defer runtime.Close()

	assert.NoError(t, runtime.RunScript("", `console.debug('foo')
console.log('bar')
console.info('baz')
console.warn('qux')
console.error('quux')`))
}

func TestRuntime_AssertEmpty(t *testing.T) {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		runtime, _  = js.NewRuntime(ctx, log.NewNop())
	)

	defer func() { cancel(); runtime.Close() }()

	for name, testCase := range map[string]struct {
		giveScript string
		wantEvents bool
		wantError  bool
	}{
		"assert.true bool":           {`assert.true(true)`, false, false},
		"assert.true bool (error)":   {`assert.true(false)`, true, false},
		"assert.true array (error)":  {`assert.true([])`, true, false},
		"mustBe.true bool":           {`mustBe.true(true)`, false, false},
		"mustBe.true bool (error)":   {`mustBe.true(false)`, true, true},
		"mustBe.true object (error)": {`mustBe.true({})`, true, true},

		"assert.false bool":           {`assert.false(false)`, false, false},
		"assert.false bool (error)":   {`assert.false(true)`, true, false},
		"assert.false array (error)":  {`assert.false([])`, true, false},
		"mustBe.false bool":           {`mustBe.false(false)`, false, false},
		"mustBe.false bool (error)":   {`mustBe.false(true)`, true, true},
		"mustBe.false object (error)": {`mustBe.false({})`, true, true},

		"assert.equals str":         {`assert.equals("foo", "foo")`, false, false},
		"assert.equals bool":        {`assert.equals(true, true)`, false, false},
		"assert.equals array":       {`assert.equals([1], [1])`, false, false},
		"assert.equals array2":      {`assert.equals([1, "foo"], [1, "foo"])`, false, false},
		"assert.equals null":        {`assert.equals(null, null)`, false, false},
		"assert.equals object":      {`assert.equals({foo: 1}, {foo: 1})`, false, false},
		"assert.equals str (error)": {`assert.equals("foo", "bar")`, true, false},
		"mustBe.equals str (error)": {`mustBe.equals("foo", "bar")`, true, true},

		"assert.notEquals str":           {`assert.notEquals("foo", "bar")`, false, false},
		"assert.notEquals bool":          {`assert.notEquals(true, false)`, false, false},
		"assert.notEquals array":         {`assert.notEquals([1], [1, 1])`, false, false},
		"assert.notEquals object":        {`assert.notEquals({foo: {bar: 1}}, {foo: {bar: 2}})`, false, false},
		"mustBe.notEquals str":           {`mustBe.notEquals("foo", "bar")`, false, false},
		"assert.notEquals str (error)":   {`assert.notEquals("foo", "foo")`, true, false},
		"mustBe.notEquals array (error)": {`mustBe.notEquals([0], [0])`, true, true},

		"assert.empty array":  {`assert.empty([])`, false, false},
		"assert.empty object": {`assert.empty({})`, false, false},
		"assert.empty bool":   {`assert.empty(false)`, false, false},
		"assert.empty null":   {`assert.empty(null)`, false, false},
		"assert.empty undef":  {`assert.empty(undefined)`, false, false},
		"assert.empty str":    {`assert.empty("")`, false, false},
		"assert.empty 0":      {`assert.empty(0)`, false, false},
		"assert.empty NaN":    {`assert.empty(NaN)`, false, false},
		"assert.empty map":    {`assert.empty(new Map())`, false, false},
		"assert.empty set":    {`assert.empty(new Set())`, false, false},
		"assert.empty date":   {`assert.empty(new Date('invalid'))`, false, false},

		"assert.empty array (error)":   {`assert.empty([1])`, true, false},
		"assert.empty object (error)":  {`assert.empty({foo: 'bar'})`, true, false},
		"assert.empty bool (error)":    {`assert.empty(true)`, true, false},
		"assert.empty array2 (error)":  {`assert.empty([undefined])`, true, false},
		"assert.empty str (error)":     {`assert.empty("str")`, true, false},
		"assert.empty 1 (error)":       {`assert.empty(1)`, true, false},
		"assert.empty symbol (error)":  {`assert.empty(Symbol())`, true, false},
		"assert.empty func (error)":    {`assert.empty(() => {})`, true, false},
		"assert.empty date (error)":    {`assert.empty(new Date())`, true, false},
		"assert.empty promise (error)": {`assert.empty(new Promise(() => {}))`, true, false},

		"mustBe.empty array":         {`mustBe.empty([])`, false, false},
		"mustBe.empty array (error)": {`mustBe.empty([1])`, true, true},

		"assert.notEmpty array":  {`assert.notEmpty([1])`, false, false},
		"assert.notEmpty object": {`assert.notEmpty({foo: 1})`, false, false},

		"assert.contains str":   {`assert.contains("foo", "foobar")`, false, false},
		"assert.contains array": {`assert.contains("foo", ["baz", "bar", "foo"])`, false, false},

		"assert.notContains str":   {`assert.notContains("foo", "bazbar")`, false, false},
		"assert.notContains array": {`assert.notContains("foo", ["baz", "bar", "yeah"])`, false, false},
		"mustBe.notContains str":   {`mustBe.notContains("foo", "foobar")`, true, true},
	} {
		tt := testCase

		t.Run(name, func(t *testing.T) {
			err := runtime.RunScript("", tt.giveScript)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			var eventsCount = len(runtime.Events())

			if tt.wantEvents {
				assert.True(t, eventsCount > 0)

				for ch := runtime.Events(); len(ch) > 0; {
					<-ch // empty the channel
				}
			} else {
				assert.Zero(t, eventsCount)
			}
		})
	}
}
