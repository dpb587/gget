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

var (
	appName                        = "gget"
	appOrigin                      = "github.com/dpb587/gget"
	appSemver, appCommit, appBuilt string
)

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
	if err != nil {
		fatal(err)
	} else if cmd.Runtime.Help {
		helpBuf := &bytes.Buffer{}
		parser.WriteHelp(helpBuf)
		help := helpBuf.String()

		// imply origin is required (optional to support --version, -h)
		help = strings.Replace(help, "[HOST/OWNER/REPOSITORY[@REF]]", "HOST/OWNER/REPOSITORY[@REF]", -1)

		// join conventional paren groups
		help = strings.Replace(help, ") (", "; ", -1)

		fmt.Print(help)
		fmt.Printf("\n")

		return
	} else if cmd.Runtime.Version != nil {
		if cmd.Runtime.Version.IsLatest {
			err := func() error {
				ref, err := service.ParseRefString(appOrigin)
				if err != nil {
					return errors.Wrap(err, "parsing app origin")
				}

				svc := github.NewService(cmd.Runtime.Logger(), github.NewClientFactory(cmd.Runtime.Logger(), cmd.Runtime.NewHTTPClient))
				res, err := svc.ResolveRef(context.Background(), service.LookupRef{Ref: ref})
				if err != nil {
					return errors.Wrap(err, "resolving app origin")
				}

				err = cmd.Runtime.Version.UnmarshalFlag(res.CanonicalRef().Ref)
				if err != nil {
					return errors.Wrap(err, "parsing version constraint")
				}

				return nil
			}()
			if err != nil {
				fatal(errors.Wrap(err, "checking latest app version"))
			}
		}

		ver, err := semver.NewVersion(v.Semver)
		if err != nil {
			fatal(errors.Wrap(err, "parsing application version"))
		}

		if !cmd.Runtime.Quiet {
			app.WriteVersion(os.Stdout, os.Args[0], v, len(cmd.Runtime.Verbose))
		}

		if !cmd.Runtime.Version.Constraint.Check(ver) {
			if cmd.Runtime.Quiet {
				os.Exit(1)
			}

			fatal(errors.Wrapf(fmt.Errorf("constraint not met: %s", cmd.Runtime.Version.Constraint.RawValue), "verifying application version"))
		}

		return
	}

	if err = cmd.Execute(nil); err != nil {
		fatal(err)
	}
}
