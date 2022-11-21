package run

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/urfave/cli/v2"

	"github.com/tarampampam/poke/internal/js"
	"github.com/tarampampam/poke/internal/js/events"
	"github.com/tarampampam/poke/internal/js/printer"
	"github.com/tarampampam/poke/internal/log"
)

type command struct {
	c *cli.Command
}

// NewCommand creates `run` command.
func NewCommand(l log.Logger) *cli.Command { //nolint:funlen
	const (
		syncFlagName              = "sync"
		threadsCountFlagName      = "threads"
		maxScriptExecTimeFlagName = "max-script-exec-time"
	)

	var cmd = command{}

	cmd.c = &cli.Command{
		Name:        "run",
		ArgsUsage:   "<files-or-directories...>",
		Aliases:     []string{"r"},
		Usage:       "Run poke files",
		Description: "Wildcards are supported, e.g. './tests/**/*.js'",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    syncFlagName,
				Aliases: []string{"s"},
				Usage:   "forces scripts running in sync mode",
			},
			&cli.UintFlag{
				Name:    threadsCountFlagName,
				Aliases: []string{"t"},
				Usage:   "number of threads to run scripts in parallel",
				Value:   uint(runtime.NumCPU() * 3), //nolint:gomnd // default value
			},
			&cli.DurationFlag{
				Name:  maxScriptExecTimeFlagName,
				Usage: "maximum execution time of each script, e.g. '10s' or '1m'",
				Value: 60 * time.Second, //nolint:gomnd // default value
			},
		},
		Action: func(c *cli.Context) error {
			var (
				threadsCount      = c.Uint(threadsCountFlagName)
				maxScriptExecTime = c.Duration(maxScriptExecTimeFlagName)
			)

			if c.Bool(syncFlagName) || threadsCount == 0 {
				threadsCount = 1
			}

			var ctx, cancel = context.WithCancel(c.Context) // main context creation
			defer cancel()

			cmd.subscribeForSystemSignals(ctx, func(_ os.Signal) { cancel() })

			files, findingErr := cmd.FindFiles(c.Args().Slice())
			if findingErr != nil {
				return findingErr
			}

			if len(files) == 0 {
				return fmt.Errorf("no files found in %s", strings.Join(c.Args().Slice(), ", "))
			}

			l.Debug("Found files", log.With("files", files))

			var (
				wg           sync.WaitGroup
				guard        = make(chan struct{}, threadsCount)
				stats        = NewOverallRunningStats()
				groupStartAt = time.Now()
			)

		runLoop:
			for _, filePath := range files {
				select {
				case guard <- struct{}{}: // would block if guard channel is already filled
					wg.Add(1)

				case <-ctx.Done():
					break runLoop
				}

				go func(filePath string) {
					defer func() { <-guard; /* release the guard */ wg.Done() }()

					startedAt := time.Now()

					l.Info("Running script", log.With("file", filePath))

					ev, runningErr := cmd.RunScript(ctx, filePath, maxScriptExecTime)

					stats.SetDuration(filePath, time.Since(startedAt))

					if runningErr != nil {
						stats.SetError(filePath, runningErr)
						l.Error("Script execution failed", log.With("file", filePath), log.With("error", runningErr))

						return
					}

					stats.SetEvents(filePath, ev)
					l.Success("Script executed successfully", log.With("file", filePath))
				}(filePath)
			}

			wg.Wait()
			close(guard)

			stats.SetSummaryDuration(time.Since(groupStartAt))

			if _, err := fmt.Fprintf(os.Stdout, "\n%s\n", stats.ToConsole()); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd.c
}

func (cmd *command) subscribeForSystemSignals(ctx context.Context, fn func(sig os.Signal)) {
	var sigs = make(chan os.Signal, 1)

	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer func() { signal.Stop(sigs); close(sigs) }()

		select {
		case <-ctx.Done():
			break

		case sig, opened := <-sigs:
			if ctx.Err() != nil { // additional context checking
				break
			}

			if opened && sig != nil {
				fn(sig)
			}
		}
	}()
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

var colorLogPrefix = text.Colors{text.FgWhite} //nolint:gochecknoglobals

func (cmd *command) RunScript( //nolint:funlen
	pCtx context.Context,
	filePath string,
	maxExecTime time.Duration,
) ([]events.Event, error) {
	script, readErr := os.ReadFile(filePath)
	if readErr != nil {
		return nil, readErr
	}

	ctx, cancel := context.WithTimeout(pCtx, maxExecTime)
	defer cancel()

	interpreter, createErr := js.NewRuntime(ctx, js.WithPrinter(printer.StringPrefixPrinter(
		colorLogPrefix.Sprintf("%s: ", filePath),
	)))

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

			case event, ok := <-interpreter.Events():
				if !ok {
					return
				}

				buf = append(buf, event)

				if event.Level == events.LevelError {
					cancel()

					interpreter.Interrupt("error event received")
				}
			}
		}
	}()

	t := time.AfterFunc(maxExecTime, func() {
		cancel()

		interpreter.Interrupt(fmt.Sprintf("script execution time exceeded (%s)", maxExecTime))
	})

	defer t.Stop()

	runErr := interpreter.RunScript(filePath, string(script))
	interpreter.Close()

	<-locker

	if runErr != nil {
		return nil, runErr
	}

	return buf, nil
}
