package run

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/urfave/cli/v2"

	"github.com/tarampampam/poke/internal/interpreter"
	"github.com/tarampampam/poke/internal/interpreter/addons"
	"github.com/tarampampam/poke/internal/syncmap"
)

// NewCommand creates `run` command.
func NewCommand() *cli.Command {
	return &cli.Command{
		Name:        "run",
		ArgsUsage:   "<files-or-directories...>",
		Aliases:     []string{"r"},
		Usage:       "Run poke files",
		Description: "Wildcards are supported, e.g. `./tests/*/*.js`",
		Action: func(c *cli.Context) error {
			var files []string

			for _, arg := range c.Args().Slice() {
				matches, err := filepath.Glob(arg)
				if err != nil {
					return err
				}

				files = append(files, matches...)
			}

			if len(files) == 0 {
				return errors.New("no files or directories provided")
			}

			var (
				wg      sync.WaitGroup
				results syncmap.SyncMap[string, []addons.Report]
			)

			for _, file := range files {
				wg.Add(1)

				go func(file string) {
					defer wg.Done()

					var runtime, err = interpreter.NewRuntime(c.Context)
					if err != nil {
						results.Store(file, []addons.Report{{
							ReportLevel: addons.ReportLevelError,
							Message:     err.Error(),
						}})

						return
					}

					script, err := os.ReadFile(file)
					if err != nil {
						results.Store(file, []addons.Report{{
							ReportLevel: addons.ReportLevelError,
							Message:     err.Error(),
						}})

						return
					}

					reports, err := runtime.RunString(string(script))
					if err != nil {
						results.Store(file, []addons.Report{{
							ReportLevel: addons.ReportLevelError,
							Message:     err.Error(),
						}})

						return
					}

					results.Store(file, reports)
				}(file)
			}

			wg.Wait()

			var hasErrors bool

			results.Range(func(file string, reports []addons.Report) bool {
				fmt.Printf("%s:", file)

				if len(reports) == 0 {
					fmt.Print(" no reports\n")
				} else {
					fmt.Print("\n")

					for _, report := range reports {
						if report.ReportLevel == addons.ReportLevelError {
							hasErrors = true
						}

						var color = text.FgBlue

						switch report.ReportLevel {
						case addons.ReportLevelError:
							color = text.FgRed

						case addons.ReportLevelWarn:
							color = text.FgYellow

						case addons.ReportLevelInfo:
							color = text.FgBlue

						case addons.ReportLevelDebug:
							color = text.FgCyan
						}

						_, _ = fmt.Fprintf(os.Stdout,
							"\t%s\n",
							text.Colors{color, text.Bold}.Sprintf("%s", report.Message),
						)
					}
				}

				return true
			})

			if hasErrors {
				return errors.New("some errors occurred")
			}

			return nil
		},
	}
}
