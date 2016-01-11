package cli

import (
	"github.com/codegangsta/cli"
)

func proxy(c *cli.Context) {

	go join(c)
	manage(c)
}
