package github

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bgentry/go-netrc/netrc"
	"github.com/dpb587/gget/pkg/service"
	"github.com/google/go-github/v29/github"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type ClientFactory struct{}

func (cf ClientFactory) Get(ctx context.Context, ref service.Ref) (*github.Client, error) {
	var tc *http.Client

	if v := os.Getenv("GITHUB_TOKEN"); v != "" {
		tc = oauth2.NewClient(
			ctx,
			oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: v},
			),
		)
	} else {
		var err error

		tc, err = cf.loadNetrc(ctx, ref)
		if err != nil {
			return nil, errors.Wrap(err, "loading auth from netrc")
		}
	}

	return github.NewClient(tc), nil
}

func (cf ClientFactory) loadNetrc(ctx context.Context, ref service.Ref) (*http.Client, error) {
	netrcPath := os.Getenv("NETRC")
	if netrcPath == "" {
		var err error

		netrcPath, err = homedir.Expand(filepath.Join("~", ".netrc"))
		if err != nil {
			return nil, errors.Wrap(err, "expanding $HOME")
		}
	}

	fi, err := os.Stat(netrcPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "checking file")
	} else if fi.IsDir() {
		// weird
		return nil, nil
	}

	rc, err := netrc.ParseFile(netrcPath)
	if err != nil {
		return nil, errors.Wrap(err, "parsing netrc")
	}

	machine := rc.FindMachine(ref.Server)
	if machine == nil {
		return nil, nil
	}

	res := oauth2.NewClient(
		ctx,
		oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: machine.Password},
		),
	)

	return res, nil
}
