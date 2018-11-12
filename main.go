package main

import (
	"os"

	// "github.com/lwyj123/storysman-ex/errs"
	"github.com/lwyj123/storysman-ex/game/cli"
)

func main() {
	if err := cli.RootCmd.Execute(); err != nil {
		// os.Exit(errs.ErrGeneral)
		os.Exit(1)
	}
}
