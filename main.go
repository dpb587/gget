package main

import (
	"os"

	"github.com/dpb587/ghet/cmd/ghet"
	"github.com/jessevdk/go-flags"
)

func main() {
	command := ghet.NewCommand()

	parser := flags.NewParser(command, flags.Default)

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	err := command.Execute(nil)
	if err != nil {
		panic(err)
	}
}
