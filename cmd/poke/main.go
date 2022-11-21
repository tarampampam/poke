// Main CLI application entrypoint.
package main

import (
	"os"

	"github.com/tarampampam/poke/internal/cli"
	"github.com/tarampampam/poke/internal/log"
)

// exitFn is a function for application exiting.
var exitFn = os.Exit //nolint:gochecknoglobals

// main CLI application entrypoint.
func main() {
	var l = log.New(log.InfoLevel)

	code, err := run(l)
	if err != nil {
		l.Fatal(err.Error())
	}

	exitFn(code)
}

// run this CLI application.
// Exit codes documentation: <https://tldp.org/LDP/abs/html/exitcodes.html>
func run(l *log.Log) (int, error) {
	if err := (cli.NewApp(l)).Run(os.Args); err != nil {
		return 1, err
	}

	return 0, nil
}
