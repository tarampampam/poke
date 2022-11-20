package cli_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tarampampam/poke/internal/cli"
)

func TestNewApp(t *testing.T) {
	app := cli.NewApp()

	require.NotEmpty(t, app.Commands)
}
