package addons

import (
	"math/rand"
	"time"

	js "github.com/dop251/goja"
	"github.com/go-faker/faker/v4"
)

type Faker struct {
	runtime *js.Runtime
	rnd     *rand.Rand
}

const (
	characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func NewFaker(runtime *js.Runtime) *Faker {
	return &Faker{
		runtime: runtime,
		rnd:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (f Faker) try(fn func() error) {
	if err := fn(); err != nil {
		panic(f.runtime.ToValue(err.Error()))
	}
}

// Bool returns a random boolean value.
func (f Faker) Bool() (b bool) {
	f.try(func() error {
		return faker.FakeData(&b)
	})

	return
}

// Falsy returns a random falsy value.
func (f Faker) Falsy() js.Value {
	switch f.rnd.Intn(6) {
	case 0:
		return f.runtime.ToValue(false)
	case 1:
		return f.runtime.ToValue(0)
	case 2:
		return f.runtime.ToValue("")
	case 3:
		return js.Null()
	case 4:
		return js.Undefined()
	default:
		return js.NaN()
	}
}

// Character returns a random character.
func (f Faker) Character(call js.FunctionCall) js.Value {
	if len(call.Arguments) == 1 {
		if poolOption := call.Argument(0).ToObject(f.runtime).Get("pool"); poolOption != nil {
			var pool = []rune(poolOption.String())

			return f.runtime.ToValue(string(pool[f.rnd.Intn(len(pool))]))
		}
	}

	return f.runtime.ToValue(string(characters[f.rnd.Intn(len(characters))]))
}

// Floating returns a random floating point number.
func (f Faker) Floating() float32 {
	return 0
}

func (f Faker) Register(runtime *js.Runtime) error {
	return runtime.GlobalObject().DefineDataProperty(
		"faker",
		runtime.ToValue(f),
		js.FLAG_FALSE, // writable
		js.FLAG_FALSE, // configurable
		js.FLAG_TRUE,  // enumerable
	)
}
