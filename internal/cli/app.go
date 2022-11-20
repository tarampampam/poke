// Package cli contains CLI command handlers.
package cli

import (
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/urfave/cli/v2"

	"github.com/tarampampam/poke/internal/cli/run"
	"github.com/tarampampam/poke/internal/env"
	"github.com/tarampampam/poke/internal/version"
)

// NewApp creates new console application.
func NewApp() *cli.App {
	return &cli.App{
		Usage: "Poke files runner",
		Before: func(context *cli.Context) error {
			if _, exists := env.ForceColors.Lookup(); exists {
				text.EnableColors()
			} else if _, exists = env.NoColors.Lookup(); exists {
				text.DisableColors()
			} else if v, ok := env.Term.Lookup(); ok && v == "dumb" {
				text.DisableColors()
			}

			return nil
		},
		Version: version.Version(),
		Commands: []*cli.Command{
			run.NewCommand(),
		},
	}
}
