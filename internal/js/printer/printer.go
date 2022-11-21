package printer

import (
	"fmt"
	"io"
	"strings"
)

// Printer is a function that prints something to somewhere :)
type Printer func(io.Writer, ...any) error

// DefaultPrinter returns default printing function.
func DefaultPrinter() Printer {
	return func(w io.Writer, args ...any) (err error) {
		_, err = fmt.Fprint(w, args...)

		return
	}
}

func StringPrefixPrinter(prefix string) Printer {
	return func(w io.Writer, args ...any) (err error) {
		const nl = "\n"

		var (
			split   = fmt.Sprint(args...)
			nlCount = strings.Count(split, nl)
			str     = prefix + strings.Replace(split, nl, nl+prefix, nlCount-1) // replace all except last
		)

		_, err = fmt.Fprint(w, str)

		return
	}
}
