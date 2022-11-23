package addons

import (
	"math/rand"
	"strings"
	"time"

	js "github.com/dop251/goja"
	"github.com/go-faker/faker/v4"
)

type Faker struct {
	runtime *js.Runtime
	rnd     *rand.Rand
}

const (
	lettersLower = "abcdefghijklmnopqrstuvwxyz"
	lettersUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	characters   = lettersLower + lettersUpper + "0123456789"
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
func (f Faker) Character(args ...js.Value) string {
	if len(args) > 0 {
		if poolOption := args[0].ToObject(f.runtime).Get("pool"); poolOption != nil {
			var pool = []rune(poolOption.String())

			return string(pool[f.rnd.Intn(len(pool))])
		}
	}

	return string(characters[f.rnd.Intn(len(characters))])
}

// Floating returns a random floating point number.
func (f Faker) Floating() (float float32) {
	f.try(func() error {
		return faker.FakeData(&float)
	})

	return
}

// Integer returns a random integer number.
func (f Faker) Integer(args ...js.Value) int {
	var min, max = -9007199254740991, 9007199254740991

	if len(args) > 0 {
		var options = args[0].ToObject(f.runtime)

		if optMin := options.Get("min"); optMin != nil {
			min = int(optMin.ToInteger())
		}

		if optMax := options.Get("max"); optMax != nil {
			max = int(optMax.ToInteger())
		}
	}

	return f.rnd.Intn(max-min) + min
}

// Letter returns a random letter.
func (f Faker) Letter() string {
	return string(lettersLower[f.rnd.Intn(len(lettersLower))])
}

// String returns a random string.
func (f Faker) String(args ...js.Value) string {
	var (
		length = 11
		pool   = []rune(characters)
	)

	if len(args) > 0 {
		var options = args[0].ToObject(f.runtime)

		if optLength := options.Get("length"); optLength != nil {
			if l := int(optLength.ToInteger()); l > 0 {
				length = l
			}
		}

		if optPool := options.Get("pool"); optPool != nil {
			pool = []rune(optPool.String())
		}
	}

	var b strings.Builder

	b.Grow(len(pool))

	for i := 0; i < length; i++ {
		b.WriteRune(pool[f.rnd.Intn(len(pool))])
	}

	return b.String()
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
