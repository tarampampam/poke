// Package cli contains CLI command handlers.
package cli

import (
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/urfave/cli/v2"

	"github.com/tarampampam/poke/internal/cli/dts"
	"github.com/tarampampam/poke/internal/cli/run"
	"github.com/tarampampam/poke/internal/env"
	"github.com/tarampampam/poke/internal/js"
	"github.com/tarampampam/poke/internal/log"
	"github.com/tarampampam/poke/internal/version"
)

// NewApp creates new console application.
func NewApp(l interface {
	log.Logger
	log.Leveler
}) *cli.App {
	const (
		logLevelFlagName = "log-level"
		verboseFlagName  = "verbose"
		defaultLogLevel  = log.InfoLevel
	)

	return &cli.App{
		Usage: "Poke files runner",
		Before: func(c *cli.Context) error {
			if _, exists := env.ForceColors.Lookup(); exists {
				text.EnableColors()
			} else if _, exists = env.NoColors.Lookup(); exists {
				text.DisableColors()
			} else if v, ok := env.Term.Lookup(); ok && v == "dumb" {
				text.DisableColors()
			}

			// parse logging level
			if logLevel, err := log.ParseLevel([]byte(c.String(logLevelFlagName))); err != nil {
				return err
			} else {
				l.SetLevel(logLevel)
			}

			if c.Bool(verboseFlagName) {
				l.SetLevel(log.DebugLevel)
			}

			return nil
		},
		Version: version.Version(),
		Commands: []*cli.Command{
			run.NewCommand(l),
			dts.NewCommand(js.DTS()),
		},
		Flags: []cli.Flag{ // global flags
			&cli.StringFlag{
				Name:    logLevelFlagName,
				Value:   defaultLogLevel.String(),
				Usage:   "logging level (" + strings.Join(log.AllLevelStrings(), "|") + ")",
				EnvVars: []string{env.LogLevel.String()},
			},
			&cli.BoolFlag{
				Name:  verboseFlagName,
				Usage: "verbose output (set the logging level to debug)",
			},
		},
	}
}
