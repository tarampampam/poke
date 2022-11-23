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
	hashLetters  = "0123456789abcdef"
)

var (
	tld = [...]string{"com", "org", "edu", "gov", "uk", "net", "io"}
)

func NewFaker(runtime *js.Runtime) *Faker {
	return &Faker{
		runtime: runtime,
		rnd:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Bool returns a random boolean value.
func (f Faker) Bool() bool {
	return f.rnd.Intn(2) == 1
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
func (f Faker) Floating() float32 {
	const min, max float32 = -32768, 32768

	return min + f.rnd.Float32()*(max-min)
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

// Paragraph returns a random paragraph.
func (f Faker) Paragraph() string { return faker.Paragraph() }

// Word returns a random word.
func (f Faker) Word() string { return faker.Word() }

// Domain returns a random domain.
func (f Faker) Domain() string { return faker.DomainName() }

// Email returns a random email.
func (f Faker) Email() string { return faker.Email() }

// Ip returns a random IPv4 address.
func (f Faker) Ip() string { return faker.IPv4() }

// Ipv6 returns a random IPv6 address.
func (f Faker) Ipv6() string { return faker.IPv6() }

// Tld returns a random top-level domain.
func (f Faker) Tld() string { return tld[f.rnd.Intn(len(tld))] }

// Url returns a random URL.
func (f Faker) Url() string { return faker.URL() }

// Date returns a random date.
func (f Faker) Date() js.Value {
	var min, max = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano(), time.Now().UnixNano()

	date, err := f.runtime.New(
		f.runtime.Get("Date").ToObject(f.runtime),
		f.runtime.ToValue((f.rnd.Int63n(max-min)+min)/1e6), // https://bit.ly/3i8c9o9
	)
	if err != nil {
		panic(f.runtime.ToValue(err.Error()))
	}

	return date
}

// Hash returns a random hash.
func (f Faker) Hash(args ...js.Value) string {
	var (
		length = 40
		pool   = []rune(hashLetters)
	)

	if len(args) > 0 {
		if optLength := args[0].ToObject(f.runtime).Get("length"); optLength != nil {
			if l := int(optLength.ToInteger()); l > 0 {
				length = l
			}
		}
	}

	var b strings.Builder

	b.Grow(len(hashLetters))

	for i := 0; i < length; i++ {
		b.WriteRune(pool[f.rnd.Intn(len(pool))])
	}

	return b.String()
}

// Uuid returns a random UUID.
func (f Faker) Uuid() string { return faker.UUIDHyphenated() }

// Random returns a randomly picked argument.
func (f Faker) Random(args ...js.Value) js.Value {
	if len(args) > 0 {
		return args[f.rnd.Intn(len(args))]
	}

	return js.Undefined()
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
