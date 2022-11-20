package run

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/urfave/cli/v2"

	"github.com/tarampampam/poke/internal/js"
	"github.com/tarampampam/poke/internal/js/events"
	"github.com/tarampampam/poke/internal/syncmap"
)

// NewCommand creates `run` command.
func NewCommand() *cli.Command { //nolint:funlen,gocognit,gocyclo
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
				wg            sync.WaitGroup
				summaryEvents syncmap.SyncMap[string, []events.Event]
			)

			for _, file := range files {
				wg.Add(1)

				go func(filePath string) {
					defer wg.Done()

					var runtime, err = js.NewRuntime(c.Context)
					if err != nil {
						summaryEvents.Store(filePath, []events.Event{{Level: events.LevelError, Error: err}})

						return
					}

					defer runtime.Close()

					wg.Add(1)

					go func() {
						var buf = make([]events.Event, 0)

						defer func() {
							if current, ok := summaryEvents.Load(filePath); ok {
								buf = append(current, buf...)
							}

							summaryEvents.Store(filePath, buf)

							wg.Done()
						}()

						for {
							select {
							case <-c.Context.Done():
								return

							case event, ok := <-runtime.Events():
								if !ok {
									return
								}

								buf = append(buf, event)
							}
						}
					}()

					var script []byte

					if script, err = os.ReadFile(filePath); err != nil {
						summaryEvents.Store(filePath, []events.Event{{Level: events.LevelError, Error: err}})

						return
					}

					if err = runtime.RunScript(filePath, string(script)); err != nil {
						summaryEvents.Store(filePath, []events.Event{{Level: events.LevelError, Error: err}})

						return
					}
				}(file)
			}

			wg.Wait()

			var hasErrors bool

			summaryEvents.Range(func(filePath string, eventsSlice []events.Event) bool {
				fmt.Printf("%s:", filePath) //nolint:forbidigo

				if len(eventsSlice) == 0 {
					fmt.Print(" no reports\n") //nolint:forbidigo
				} else {
					fmt.Print("\n") //nolint:forbidigo

					for _, e := range eventsSlice {
						if e.Level == events.LevelError {
							hasErrors = true
						}

						var color = text.FgBlue

						switch e.Level {
						case events.LevelError:
							color = text.FgRed

						case events.LevelWarn:
							color = text.FgYellow

						case events.LevelInfo:
							color = text.FgBlue

						case events.LevelDebug:
							color = text.FgCyan
						}

						_, _ = fmt.Fprintf(os.Stdout,
							"\t%s\n",
							text.Colors{color, text.Bold}.Sprintf("%s (%s)", e.Message, e.Error),
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
