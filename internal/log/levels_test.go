package log_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tarampampam/poke/internal/log"
)

func TestLevel_String(t *testing.T) {
	for name, tt := range map[string]struct {
		giveLevel  log.Level
		wantString string
	}{
		"debug":     {giveLevel: log.DebugLevel, wantString: "debug"},
		"info":      {giveLevel: log.InfoLevel, wantString: "info"},
		"warn":      {giveLevel: log.WarnLevel, wantString: "warn"},
		"error":     {giveLevel: log.ErrorLevel, wantString: "error"},
		"<unknown>": {giveLevel: log.Level(127), wantString: "level(127)"},
	} {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.wantString, tt.giveLevel.String())
		})
	}
}

func TestParseLevel(t *testing.T) {
	for name, tt := range map[string]struct {
		giveText  []byte
		wantLevel log.Level
		wantError error
	}{
		"<empty value>": {giveText: []byte(""), wantLevel: log.InfoLevel},
		"trace":         {giveText: []byte("debug"), wantLevel: log.DebugLevel},
		"verbose":       {giveText: []byte("debug"), wantLevel: log.DebugLevel},
		"debug":         {giveText: []byte("debug"), wantLevel: log.DebugLevel},
		"info":          {giveText: []byte("info"), wantLevel: log.InfoLevel},
		"warn":          {giveText: []byte("warn"), wantLevel: log.WarnLevel},
		"error":         {giveText: []byte("error"), wantLevel: log.ErrorLevel},
		"foobar":        {giveText: []byte("foobar"), wantError: errors.New("unrecognized logging level: \"foobar\"")},
	} {
		t.Run(name, func(t *testing.T) {
			l, err := log.ParseLevel(tt.giveText)

			if tt.wantError == nil {
				require.NoError(t, err)
				require.Equal(t, tt.wantLevel, l)
			} else {
				require.EqualError(t, err, tt.wantError.Error())
			}
		})
	}
}
