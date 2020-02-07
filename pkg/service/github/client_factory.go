package github

import (
	"context"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/bgentry/go-netrc/netrc"
	"github.com/dpb587/gget/pkg/service"
	"github.com/google/go-github/v29/github"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type ClientFactory struct {
	RoundTripFactory func(http.RoundTripper) http.RoundTripper
}

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

	if tc == nil {
		tc = &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).Dial,
				IdleConnTimeout:       15 * time.Second,
				TLSHandshakeTimeout:   15 * time.Second,
				ResponseHeaderTimeout: 15 * time.Second,
				ExpectContinueTimeout: 5 * time.Second,
			},
		}
	}

	if cf.RoundTripFactory != nil {
		tc.Transport = cf.RoundTripFactory(tc.Transport)
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
