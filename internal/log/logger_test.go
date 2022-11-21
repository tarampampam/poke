package log_test

import (
	"bytes"
	"regexp"
	"sync"
	"testing"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"

	"github.com/tarampampam/poke/internal/log"
)

func TestFoo(t *testing.T) { // TODO delete me
	var (
		l     = log.New(log.DebugLevel)
		extra = []log.Extra{
			log.With("string", "foo"),
			log.With("int", 123),
			log.With("struct", struct{}{}),
			log.With("slice", []string{"bar"}),
		}
	)

	l.Debug("debug msg")
	l.Info("", extra...)
	l.Info("info \nmsginfo \nmsginfo \nmsginfo \nmsginfo \nmsg", extra...)
	l.Success("success msg", extra...)
	l.Warn("warn msg", extra...)
	l.Error("error msg", extra...)
	l.Fatal("fatal msg", extra...)
}

func TestLog_Debug(t *testing.T) {
	text.DisableColors()
	defer text.EnableColors()

	var (
		stdOut, errOut bytes.Buffer
		extra          = []log.Extra{
			log.With("string", "foo"),
			log.With("int", 123),
			log.With("struct", struct{}{}),
			log.With("slice", []string{"bar"}),
		}
	)

	var l = log.New(log.DebugLevel, log.WithStdOut(&stdOut), log.WithStdErr(&errOut))

	l.Debug("debug msg", extra...)
	l.Info("info msg", extra...)
	l.Warn("warn msg", extra...)
	l.Error("error msg", extra...)

	var (
		stdOutStr = stdOut.String()
		stdErrStr = errOut.String()
	)

	assert.Contains(t, stdOutStr, "debug msg")
	assert.Contains(t, stdOutStr, "info msg")
	assert.Contains(t, stdOutStr, "warn msg")
	assert.NotContains(t, stdOutStr, "error msg")

	assert.NotContains(t, stdErrStr, "debug msg")
	assert.NotContains(t, stdErrStr, "info msg")
	assert.NotContains(t, stdErrStr, "warn msg")
	assert.Contains(t, stdErrStr, "error msg")

	assert.Regexp(t, regexp.MustCompile(`^\s+debug\s+debug msg`), stdOutStr)
	assert.Contains(t, stdOutStr, "string:foo")
	assert.Contains(t, stdOutStr, "int:123")
	assert.Contains(t, stdOutStr, "struct:{}")
	assert.Contains(t, stdOutStr, "slice:[bar]")
}

func TestLog_Info(t *testing.T) {
	text.DisableColors()
	defer text.EnableColors()

	var (
		stdOut, errOut bytes.Buffer
		extra          = []log.Extra{
			log.With("string", "foo"),
			log.With("int", 123),
			log.With("struct", struct{}{}),
			log.With("slice", []string{"bar"}),
		}
	)

	var l = log.New(log.InfoLevel, log.WithStdOut(&stdOut), log.WithStdErr(&errOut))

	l.Debug("debug msg", extra...)
	l.Info("info msg", extra...)
	l.Warn("warn msg", extra...)
	l.Error("error msg", extra...)

	var (
		stdOutStr = stdOut.String()
		stdErrStr = errOut.String()
	)

	assert.NotContains(t, stdOutStr, "debug msg")
	assert.Contains(t, stdOutStr, "info msg")
	assert.Contains(t, stdOutStr, "warn msg")
	assert.NotContains(t, stdOutStr, "error msg")

	assert.NotContains(t, stdErrStr, "debug msg")
	assert.NotContains(t, stdErrStr, "info msg")
	assert.NotContains(t, stdErrStr, "warn msg")
	assert.Contains(t, stdErrStr, "error msg")

	assert.Regexp(t, regexp.MustCompile(`^\s+info\s+info msg`), stdOutStr)
	assert.Contains(t, stdOutStr, "string:foo")
	assert.Contains(t, stdOutStr, "int:123")
	assert.Contains(t, stdOutStr, "struct:{}")
	assert.Contains(t, stdOutStr, "slice:[bar]")
}

func TestLog_Warn(t *testing.T) {
	text.DisableColors()
	defer text.EnableColors()

	var (
		stdOut, errOut bytes.Buffer
		extra          = []log.Extra{
			log.With("string", "foo"),
			log.With("int", 123),
			log.With("struct", struct{}{}),
			log.With("slice", []string{"bar"}),
		}
	)

	var l = log.New(log.WarnLevel, log.WithStdOut(&stdOut), log.WithStdErr(&errOut))

	l.Debug("debug msg", extra...)
	l.Info("info msg", extra...)
	l.Warn("warn msg", extra...)
	l.Error("error msg", extra...)

	var (
		stdOutStr = stdOut.String()
		stdErrStr = errOut.String()
	)

	assert.NotContains(t, stdOutStr, "debug msg")
	assert.NotContains(t, stdOutStr, "info msg")
	assert.Contains(t, stdOutStr, "warn msg")
	assert.NotContains(t, stdOutStr, "error msg")

	assert.NotContains(t, stdErrStr, "debug msg")
	assert.NotContains(t, stdErrStr, "info msg")
	assert.NotContains(t, stdErrStr, "warn msg")
	assert.Contains(t, stdErrStr, "error msg")

	assert.Regexp(t, regexp.MustCompile(`^\s+warn\s+warn msg`), stdOutStr)
	assert.Contains(t, stdOutStr, "string:foo")
	assert.Contains(t, stdOutStr, "int:123")
	assert.Contains(t, stdOutStr, "struct:{}")
	assert.Contains(t, stdOutStr, "slice:[bar]")
}

func TestLog_Error(t *testing.T) {
	text.DisableColors()
	defer text.EnableColors()

	var (
		stdOut, errOut bytes.Buffer
		extra          = []log.Extra{
			log.With("string", "foo"),
			log.With("int", 123),
			log.With("struct", struct{}{}),
			log.With("slice", []string{"bar"}),
		}
	)

	var l = log.New(log.ErrorLevel, log.WithStdOut(&stdOut), log.WithStdErr(&errOut))

	l.Debug("debug msg", extra...)
	l.Info("info msg", extra...)
	l.Warn("warn msg", extra...)
	l.Error("error msg", extra...)

	var (
		stdOutStr = stdOut.String()
		stdErrStr = errOut.String()
	)

	assert.Empty(t, stdOutStr)

	assert.NotContains(t, stdErrStr, "debug msg")
	assert.NotContains(t, stdErrStr, "info msg")
	assert.NotContains(t, stdErrStr, "warn msg")
	assert.Contains(t, stdErrStr, "error msg")

	assert.Regexp(t, regexp.MustCompile(`^\s+error\s+error msg`), stdErrStr)
	assert.Contains(t, stdErrStr, "string:foo")
	assert.Contains(t, stdErrStr, "int:123")
	assert.Contains(t, stdErrStr, "struct:{}")
	assert.Contains(t, stdErrStr, "slice:[bar]")
}

func TestLog_Concurrent(t *testing.T) {
	var (
		stdOut, errOut bytes.Buffer

		l  = log.New(log.DebugLevel, log.WithStdOut(&stdOut), log.WithStdErr(&errOut))
		wg sync.WaitGroup
	)

	for i := 0; i < 100; i++ {
		wg.Add(4)

		go func() { defer wg.Done(); l.Debug("debug", log.With("struct", struct{}{})) }()
		go func() { defer wg.Done(); l.Info("info", log.With("struct", struct{}{})) }()
		go func() { defer wg.Done(); l.Warn("warn", log.With("struct", struct{}{})) }()
		go func() { defer wg.Done(); l.Error("error", log.With("struct", struct{}{})) }()
	}

	wg.Wait()

	assert.NotEmpty(t, stdOut.String())
	assert.NotEmpty(t, errOut.String())
}
