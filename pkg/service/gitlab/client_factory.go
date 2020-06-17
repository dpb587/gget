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
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

type roundTripTransformer func(http.RoundTripper) http.RoundTripper

type ClientFactory struct {
	log              *logrus.Logger
	roundTripFactory func(http.RoundTripper) http.RoundTripper
}

func NewClientFactory(log *logrus.Logger, roundTripFactory roundTripTransformer) *ClientFactory {
	return &ClientFactory{
		log:              log,
		roundTripFactory: roundTripFactory,
	}
}

func (cf ClientFactory) Get(ctx context.Context, lookupRef service.LookupRef) (*gitlab.Client, error) {
	var tc *http.Client
	var token string

	if v := os.Getenv("GITLAB_TOKEN"); v != "" {
		cf.log.Infof("found authentication for %s: env GITLAB_TOKEN", lookupRef.Ref.Server)

		token = v
	} else {
		var err error

		token, err = cf.loadNetrc(ctx, lookupRef)
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

	if cf.roundTripFactory != nil {
		tc.Transport = cf.roundTripFactory(tc.Transport)
	}

	res, err := gitlab.NewClient(token, gitlab.WithHTTPClient(tc))
	if err != nil {
		return nil, errors.Wrap(err, "creating client")
	}

	return res, nil
}

func (cf ClientFactory) loadNetrc(ctx context.Context, lookupRef service.LookupRef) (string, error) {
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

	machine := rc.FindMachine(lookupRef.Ref.Server)
	if machine == nil {
		return "", nil
	}

	cf.log.Infof("found authentication for %s: netrc %s", lookupRef.Ref.Server, netrcPath)

	return machine.Password, nil
}
