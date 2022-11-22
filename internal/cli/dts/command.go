package dts

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

// NewCommand creates `dts` command.
func NewCommand(content string) *cli.Command {
	return &cli.Command{
		Name:  "dts",
		Usage: "Print the 'common.d.ts' file",
		Action: func(c *cli.Context) (err error) {
			_, err = fmt.Fprint(os.Stdout, content)

			return
		},
	}
}
