package gitlab

import (
	"context"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/bgentry/go-netrc/netrc"
	"github.com/dpb587/gget/pkg/service"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

type ClientFactory struct {
	RoundTripFactory func(http.RoundTripper) http.RoundTripper
}

func (cf ClientFactory) Get(ctx context.Context, ref service.Ref) (*gitlab.Client, error) {
	var tc *http.Client
	var token string

	if v := os.Getenv("GITLAB_TOKEN"); v != "" {
		token = v
	} else {
		var err error

		token, err = cf.loadNetrc(ctx, ref)
		if err != nil {
			return nil, errors.Wrap(err, "loading auth from netrc")
		}
	}

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

	if cf.RoundTripFactory != nil {
		tc.Transport = cf.RoundTripFactory(tc.Transport)
	}

	res, err := gitlab.NewClient(token, gitlab.WithHTTPClient(tc))
	if err != nil {
		return nil, errors.Wrap(err, "creating client")
	}

	return res, nil
}

func (cf ClientFactory) loadNetrc(ctx context.Context, ref service.Ref) (string, error) {
	netrcPath := os.Getenv("NETRC")
	if netrcPath == "" {
		var err error

		netrcPath, err = homedir.Expand(filepath.Join("~", ".netrc"))
		if err != nil {
			return "", errors.Wrap(err, "expanding $HOME")
		}
	}

	fi, err := os.Stat(netrcPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}

		return "", errors.Wrap(err, "checking file")
	} else if fi.IsDir() {
		// weird
		return "", nil
	}

	rc, err := netrc.ParseFile(netrcPath)
	if err != nil {
		return "", errors.Wrap(err, "parsing netrc")
	}

	machine := rc.FindMachine(ref.Server)
	if machine == nil {
		return "", nil
	}

	return machine.Password, nil
}
