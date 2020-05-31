package gget

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dpb587/gget/pkg/app"
	"github.com/dpb587/gget/pkg/service"
	"github.com/dpb587/gget/pkg/service/github"
	"github.com/dpb587/gget/pkg/service/gitlab"
	"github.com/sirupsen/logrus"
)

type Runtime struct {
	Quiet    bool   `long:"quiet" short:"q" description:"suppress runtime status reporting"`
	Verbose  []bool `long:"verbose" short:"v" description:"increase logging verbosity"`
	Parallel int    `long:"parallel" description:"maximum number of parallel operations" default:"3"`
	Service  string `long:"service" description:"specific git service to use (i.e. github, gitlab; default: auto-detect)"`

	Help    bool `long:"help" short:"h" description:"show documentation of this command"`
	Version bool `long:"version" description:"show version of this command"`

	app    app.Version
	logger *logrus.Logger
}

func NewRuntime(app app.Version) *Runtime {
	return &Runtime{
		app: app,
	}
}

func (r *Runtime) RefResolver() (service.RefResolver, error) {
	var resolvers []service.ConditionalRefResolver

	if r.Service == "" || r.Service == "github" {
		resolvers = append(
			resolvers,
			github.NewService(r.Logger(), github.NewClientFactory(r.Logger(), r.RoundTripLogger)),
		)
	}

	if r.Service == "" || r.Service == "gitlab" {
		resolvers = append(
			resolvers,
			gitlab.NewService(r.Logger(), gitlab.NewClientFactory(r.Logger(), r.RoundTripLogger)),
		)
	}

	switch len(resolvers) {
	case 0:
		return nil, fmt.Errorf("unsupported service: %s", r.Service)
	case 1:
		return resolvers[0], nil
	}

	return service.NewMultiRefResolver(resolvers...), nil
}

func (r *Runtime) Logger() *logrus.Logger {
	if r.logger == nil {
		var logLevel logrus.Level

		switch len(r.Verbose) {
		case 0:
			logLevel = logrus.FatalLevel
		case 1:
			logLevel = logrus.WarnLevel
		case 2:
			logLevel = logrus.InfoLevel
		default:
			logLevel = logrus.DebugLevel
		}

		r.logger = logrus.New()
		r.logger.Level = logLevel
		r.logger.Out = os.Stderr

		r.logger.Infof("starting %s", r.app.String())
	}

	return r.logger
}

func (r *Runtime) RoundTripLogger(rt http.RoundTripper) http.RoundTripper {
	if rt == nil {
		rt = http.DefaultTransport
	}

	return roundTripLogger{
		l:  r.logger,
		rt: rt,
	}
}

type roundTripLogger struct {
	l  *logrus.Logger
	rt http.RoundTripper
}

func (rtl roundTripLogger) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	rtl.l.Debugf("http: %s %s", req.Method, req.URL.String())

	res, err := rtl.rt.RoundTrip(req)

	rtl.l.Infof("http: %s %s (status: %s)", req.Method, req.URL.String(), res.Status)

	if v := res.Header.Get("ratelimit-remaining"); v != "" {
		rtl.l.Debugf("http: %s %s (ratelimit-remaining: %s)", req.Method, req.URL.String(), v)
	} else if v := res.Header.Get("x-ratelimit-remaining"); v != "" {
		rtl.l.Debugf("http: %s %s (x-ratelimit-remaining: %s)", req.Method, req.URL.String(), v)
	}

	return res, err
}
