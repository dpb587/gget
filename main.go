package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/dpb587/gget/cmd/gget"
	"github.com/dpb587/gget/pkg/app"
	"github.com/dpb587/gget/pkg/service"
	"github.com/dpb587/gget/pkg/service/github"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
)

var appName = "gget"
var appSemver, appCommit, appBuilt, appOrigin string

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
		buf := &bytes.Buffer{}
		parser.WriteHelp(buf)

		fmt.Print(strings.Replace(buf.String(), ") (", "; ", -1))
		fmt.Printf("\n")

		return
	} else if cmd.Runtime.Version != nil {
		if cmd.Runtime.Version.RawConstraint == "*" {
			app.WriteVersion(os.Stdout, os.Args[0], v, len(cmd.Runtime.Verbose))
		} else {
			if cmd.Runtime.Version.RawConstraint == "latest" {
				ref, err := service.ParseRefString(appOrigin)
				if err != nil {
					fatal(errors.Wrap(err, "parsing app origin"))
				}

				svc := github.NewService(cmd.Runtime.Logger(), github.NewClientFactory(cmd.Runtime.Logger(), cmd.Runtime.NewHTTPClient))
				res, err := svc.ResolveRef(context.Background(), service.LookupRef{Ref: ref})
				if err != nil {
					fatal(errors.Wrap(err, "resolving ref"))
				}

				err = cmd.Runtime.Version.UnmarshalFlag(res.CanonicalRef().Ref)
				if err != nil {
					fatal(errors.Wrap(err, "parsing constraint"))
				}
			}

			ver, err := semver.NewVersion(v.Semver)
			if err != nil {
				fatal(errors.Wrap(err, "parsing application version"))
			}

			if !cmd.Runtime.Version.Check(ver) {
				if cmd.Runtime.Quiet {
					os.Exit(1)
				}

				fatal(fmt.Errorf("version '%s' does not satisfy constraint: %s", v.Semver, cmd.Runtime.Version.RawConstraint))
			} else if !cmd.Runtime.Quiet {
				app.WriteVersion(os.Stdout, os.Args[0], v, len(cmd.Runtime.Verbose))
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
