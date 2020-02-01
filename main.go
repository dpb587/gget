package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/dpb587/ghet/cmd/ghet"
	"github.com/jessevdk/go-flags"
)

func main() {
	command := ghet.NewCommand()

	parser := flags.NewParser(command, flags.PassDoubleDash)

	fatal := func(err error) {
		if debug, _ := strconv.ParseBool(os.Getenv("DEBUG")); debug {
			panic(err)
		}

		fmt.Fprintf(os.Stderr, "%s: %s\n", parser.Command.Name, err)

		os.Exit(1)
	}

	_, err := parser.Parse()
	if command.Runtime.Help {
		parser.WriteHelp(os.Stdout)
		fmt.Printf("\n")

		return
	} else if command.Runtime.Version {
		fmt.Println("TODO")

		return
	} else if err != nil {
		fatal(err)
	}

	if err = command.Execute(nil); err != nil {
		fatal(err)
	}
}
