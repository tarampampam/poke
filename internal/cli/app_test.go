package cli_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tarampampam/poke/internal/cli"
	"github.com/tarampampam/poke/internal/log"
)

func TestNewApp(t *testing.T) {
	app := cli.NewApp(log.NewNop())

	require.NotEmpty(t, app.Commands)
}
