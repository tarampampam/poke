package run

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/urfave/cli/v2"

	"github.com/tarampampam/poke/internal/js"
	"github.com/tarampampam/poke/internal/js/events"
)

type command struct {
	c *cli.Command
}

// NewCommand creates `run` command.
func NewCommand() *cli.Command { //nolint:funlen
	var cmd = command{}

	cmd.c = &cli.Command{
		Name:        "run",
		ArgsUsage:   "<files-or-directories...>",
		Aliases:     []string{"r"},
		Usage:       "Run poke files",
		Description: "Wildcards are supported, e.g. `./tests/**/*.js`",
		Flags: []cli.Flag{
			&cli.BoolFlag{ // TODO: use it
				Name:  "sync",
				Usage: "Run scripts in sync mode",
			},
			&cli.BoolFlag{ // TODO: use it
				Name:  "async",
				Usage: "Run scripts in async mode",
			},
		},
		Action: func(c *cli.Context) error {
			files, findingErr := cmd.FindFiles(c.Args().Slice())
			if findingErr != nil {
				return findingErr
			}

			if len(files) == 0 {
				return fmt.Errorf(
					"no valid files was found in %s (check the files path and a shebang in it)",
					strings.Join(c.Args().Slice(), ", "),
				)
			}

			var (
				wg           sync.WaitGroup
				stats        = NewOverallRunningStats()
				groupStartAt = time.Now()
			)

			for _, filePath := range files {
				wg.Add(1)

				go func(filePath string) { // TODO: limit goroutines count
					defer wg.Done()

					startedAt := time.Now()

					ev, runningErr := cmd.RunScript(c.Context, filePath)

					stats.SetDuration(filePath, time.Since(startedAt))

					if runningErr != nil {
						stats.SetError(filePath, runningErr)

						return
					}

					stats.SetEvents(filePath, ev)
				}(filePath)
			}

			wg.Wait()

			stats.SetSummaryDuration(time.Since(groupStartAt))

			if _, err := fmt.Fprintf(os.Stdout, "\n%s\n", stats.ToConsole()); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd.c
}

func (cmd *command) FindFiles(in []string) ([]string, error) {
	var files []string

	for _, arg := range in {
		matches, globErr := doublestar.FilepathGlob(arg)
		if globErr != nil {
			return nil, globErr
		}

		files = append(files, matches...)
	}

	return files, nil
}

func (cmd *command) RunScript(ctx context.Context, filePath string) ([]events.Event, error) {
	script, readErr := os.ReadFile(filePath)
	if readErr != nil {
		return nil, readErr
	}

	runtime, createErr := js.NewRuntime(ctx)
	if createErr != nil {
		return nil, createErr
	}

	var (
		buf    = make([]events.Event, 0, 16) //nolint:gomnd // pre-allocation is better than re-allocation
		locker = make(chan struct{})
	)

	go func() {
		defer close(locker)

		for {
			select {
			case <-ctx.Done():
				return

			case event, ok := <-runtime.Events():
				if !ok {
					return
				}

				buf = append(buf, event)

				if event.Level == events.LevelError {
					runtime.Interrupt("error event received")
				}
			}
		}
	}()

	runErr := runtime.RunScript(filePath, string(script))
	runtime.Close()

	<-locker

	if runErr != nil {
		return nil, runErr
	}

	return buf, nil // TODO return last captured events error
}
