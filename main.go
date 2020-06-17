package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Masterminds/semver"
	"github.com/dpb587/gget/cmd/gget"
	"github.com/dpb587/gget/pkg/app"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
)

var appName = "gget"
var appSemver, appCommit, appBuilt string

func main() {
	v := app.MustVersion(appName, appSemver, appCommit, appBuilt)
	cmd := gget.NewCommand(v)

	parser := flags.NewParser(cmd, flags.PassDoubleDash)

	fatal := func(err error) {
		if debug, _ := strconv.ParseBool(os.Getenv("DEBUG")); debug {
			panic(err)
		}

		fmt.Fprintf(os.Stderr, "%s: error: %s\n", parser.Command.Name, err)

		os.Exit(1)
	}

	_, err := parser.Parse()
	if cmd.Runtime.Help {
		parser.WriteHelp(os.Stdout)
		fmt.Printf("\n")

		return
	} else if cmd.Runtime.Version != nil {
		app.WriteVersion(os.Stdout, os.Args[0], v, len(cmd.Runtime.Verbose))

		if cmd.Runtime.Version.RawConstraint != "*" {
			ver, err := semver.NewVersion(v.Semver)
			if err != nil {
				fatal(errors.Wrap(err, "parsing application version"))
			}

			if !cmd.Runtime.Version.Check(ver) {
				fatal(fmt.Errorf("version '%s' does not satisfy constraint: %s", v.Semver, cmd.Runtime.Version.RawConstraint))
			}
		}

		return
	} else if err != nil {
		fatal(err)
	}

	if err = cmd.Execute(nil); err != nil {
		fatal(err)
	}
}
