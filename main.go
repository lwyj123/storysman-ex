package main

import (
	"os"

	"github.com/lwyj123/storysman-ex/cli"
	"github.com/lwyj123/storysman-ex/errs"
)

func main() {
	if err := cli.RootCmd.Execute(); err != nil {
		os.Exit(errs.ErrGeneral)
	}
}
