package github

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bgentry/go-netrc/netrc"
	"github.com/dpb587/gget/pkg/service"
	"github.com/google/go-github/v29/github"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type roundTripTransformer func(http.RoundTripper) http.RoundTripper

type ClientFactory struct {
	log               *logrus.Logger
	httpClientFactory func() *http.Client
}

func NewClientFactory(log *logrus.Logger, httpClientFactory func() *http.Client) *ClientFactory {
	return &ClientFactory{
		log:               log,
		httpClientFactory: httpClientFactory,
	}
}

func (cf ClientFactory) Get(ctx context.Context, lookupRef service.LookupRef) (*github.Client, error) {
	var tokenSource oauth2.TokenSource

	if v := os.Getenv("GITHUB_TOKEN"); v != "" {
		cf.log.Infof("found authentication for %s: env GITHUB_TOKEN", lookupRef.Ref.Server)

		tokenSource = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: v})
	} else {
		var err error

		tokenSource, err = cf.loadNetrc(ctx, lookupRef)
		if err != nil {
			return nil, errors.Wrap(err, "loading auth from netrc")
		}
	}

	httpClient := cf.httpClientFactory()

	if tokenSource != nil {
		httpClient.Transport = &oauth2.Transport{
			Base:   httpClient.Transport,
			Source: oauth2.ReuseTokenSource(nil, tokenSource),
		}
	}

	if lookupRef.Ref.Server == "github.com" {
		return github.NewClient(httpClient), nil
	}

	// TODO figure out https configurability
	baseURL := fmt.Sprintf("https://%s", lookupRef.Ref.Server)

	c, err := github.NewEnterpriseClient(baseURL, baseURL, httpClient)
	if err != nil {
		return nil, errors.Wrap(err, "building enterprise client")
	}

	return c, nil
}

func (cf ClientFactory) loadNetrc(ctx context.Context, lookupRef service.LookupRef) (oauth2.TokenSource, error) {
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

	machine := rc.FindMachine(lookupRef.Ref.Server)
	if machine == nil {
		return nil, nil
	}

	cf.log.Infof("found authentication for %s: netrc %s", lookupRef.Ref.Server, netrcPath)

	return oauth2.StaticTokenSource(&oauth2.Token{AccessToken: machine.Password}), nil
}
