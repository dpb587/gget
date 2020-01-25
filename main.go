package main

import (
	"os"

	"github.com/dpb587/ghet/cmd/ghet"
	"github.com/jessevdk/go-flags"
)

func main() {
	main := ghet.New()

	var parser = flags.NewParser(&main, flags.Default)

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			panic(err)

			os.Exit(1)
		}
	}
}
