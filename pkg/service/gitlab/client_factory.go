package gitlab

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bgentry/go-netrc/netrc"
	"github.com/dpb587/gget/pkg/service"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
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

func (cf ClientFactory) Get(ctx context.Context, lookupRef service.LookupRef) (*gitlab.Client, error) {
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

	res, err := gitlab.NewClient(token, gitlab.WithHTTPClient(cf.httpClientFactory()))
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
