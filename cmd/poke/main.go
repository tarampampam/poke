// Main CLI application entrypoint.
package main

import (
	"fmt"
	"os"

	"github.com/tarampampam/poke/internal/playground"
)

// exitFn is a function for application exiting.
var exitFn = os.Exit //nolint:gochecknoglobals

// main CLI application entrypoint.
func main() {
	code, err := run()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}

	exitFn(code)
}

// run this CLI application.
// Exit codes documentation: <https://tldp.org/LDP/abs/html/exitcodes.html>
func run() (int, error) {
	if err := playground.Run(); err != nil {
		return 1, err
	}

	return 0, nil
}
