package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type (
	Logger interface {
		// Debug logs a message at DebugLevel.
		Debug(msg string, v ...Extra)

		// Info logs a message at InfoLevel.
		Info(msg string, v ...Extra)

		// Success logs a success message at InfoLevel.
		Success(msg string, v ...Extra)

		// Warn logs a message at WarnLevel.
		Warn(msg string, v ...Extra)

		// Error logs a message at ErrorLevel.
		Error(msg string, v ...Extra)

		// Fatal logs a message at ErrorLevel.
		Fatal(msg string, v ...Extra)
	}

	Extra interface {
		// Key returns the key of the extra field.
		Key() string

		// Value returns the value of the extra field.
		Value() any
	}
)

// extra is a helper struct for that implements Extra interface.
type extra struct {
	key   string
	value any
}

func (e *extra) Key() string { return e.key }
func (e *extra) Value() any  { return e.value }

// With returns an Extra logger field.
func With(key string, value any) Extra { return &extra{key: key, value: value} }

// Option is a function that can be used to modify a Log.
type Option func(*Log)

// WithStdOut sets the writer for the standard output.
func WithStdOut(w io.Writer) Option { return func(l *Log) { l.stdOut = w } }

// WithStdErr sets the writer for the error output.
func WithStdErr(w io.Writer) Option { return func(l *Log) { l.stdErr = w } }

// Log is a logger that logs messages at specified level.
type Log struct {
	lvl    Level
	stdOut io.Writer
	stdErr io.Writer
}

var _ Logger = (*Log)(nil) // verify that the Log implements the Logger interface

const (
	debugPrefix   = " debug "
	infoPrefix    = "  info "
	successPrefix = "    ok "
	warnPrefix    = "  warn "
	errorPrefix   = " error "
)

// NewNop creates a no-op Logger.
func NewNop() *Log {
	return &Log{
		stdOut: io.Discard,
		stdErr: io.Discard,
	}
}

// New creates a new Logger with specified level.
func New(lvl Level, opts ...Option) *Log {
	var log = &Log{
		lvl:    lvl,
		stdOut: os.Stdout,
		stdErr: os.Stderr,
	}

	for _, opt := range opts {
		opt(log)
	}

	return log
}

func (l *Log) check(lvl Level) bool {
	return lvl >= l.lvl
}

func (l *Log) write(w io.Writer, c colors, prefix, sep, msg string, extra ...Extra) {
	var b strings.Builder

	b.Grow(len(prefix) + len(msg) + 32 + len(extra)*64) // extra bytes for the ASCII color codes
	b.WriteString(c[0].Sprint(prefix))
	b.WriteString(sep)
	b.WriteString(c[1].Sprint(msg))

	if len(extra) > 0 {
		b.WriteRune('\t')

		for i, e := range extra {
			b.WriteString(c[2].Sprint(e.Key(), ":", e.Value()))

			if i < len(extra)-1 {
				b.WriteRune(' ')
			}
		}
	}

	_, _ = fmt.Fprintln(w, b.String())
}

// SetLevel sets the log level.
func (l *Log) SetLevel(lvl Level) { l.lvl = lvl }

// Debug logs a message at DebugLevel.
func (l *Log) Debug(msg string, v ...Extra) {
	if !l.check(DebugLevel) || l.stdOut == io.Discard {
		return
	}

	if _, file, line, ok := runtime.Caller(1); ok {
		v = append([]Extra{With("caller", filepath.Base(file)+":"+strconv.Itoa(line))}, v...)
	}

	l.write(l.stdOut, colorsDebug, debugPrefix, " ", msg, v...)
}

// Info logs a message at InfoLevel.
func (l *Log) Info(msg string, v ...Extra) {
	if !l.check(InfoLevel) || l.stdOut == io.Discard {
		return
	}

	l.write(l.stdOut, colorsInfo, infoPrefix, " ", msg, v...)
}

// Success logs a success message at InfoLevel.
func (l *Log) Success(msg string, v ...Extra) {
	if !l.check(InfoLevel) || l.stdOut == io.Discard {
		return
	}

	l.write(l.stdOut, colorsSuccess, successPrefix, " ", msg, v...)
}

// Warn logs a message at WarnLevel.
func (l *Log) Warn(msg string, v ...Extra) {
	if !l.check(WarnLevel) || l.stdOut == io.Discard {
		return
	}

	l.write(l.stdOut, colorsWarn, warnPrefix, " ", msg, v...)
}

// Error logs a message at ErrorLevel.
func (l *Log) Error(msg string, v ...Extra) {
	if !l.check(ErrorLevel) || l.stdErr == io.Discard {
		return
	}

	l.write(l.stdErr, colorsError, errorPrefix, " ", msg, v...)
}

// Fatal logs a message at ErrorLevel.
func (l *Log) Fatal(msg string, v ...Extra) {
	if !l.check(ErrorLevel) || l.stdErr == io.Discard {
		return
	}

	var b strings.Builder

	b.Grow(len(msg) + 4) // more faster fmt.Sprintf("  %s  ") replacement
	b.WriteString("  ")
	b.WriteString(msg)
	b.WriteString("  ")

	l.write(l.stdErr, colorsFatal, "  Fatal error  ", "", b.String(), v...)
}
