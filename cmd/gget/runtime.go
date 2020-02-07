package gget

import (
	"net/http"
	"os"

	"github.com/dpb587/gget/pkg/app"
	"github.com/sirupsen/logrus"
)

type Runtime struct {
	Quiet    bool   `long:"quiet" short:"q" description:"suppress runtime status reporting"`
	Verbose  []bool `long:"verbose" short:"v" description:"increase logging verbosity"`
	Parallel int    `long:"parallel" description:"maximum number of parallel operations" default:"3"`

	Help    bool `long:"help" short:"h" description:"show documentation of this tool"`
	Version bool `long:"version" description:"show version of this tool"`

	app    app.Version
	logger *logrus.Logger
}

func NewRuntime(app app.Version) *Runtime {
	return &Runtime{
		app: app,
	}
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
	rtl.l.Debugf("http: request: %s %s", req.Method, req.URL.String())

	res, err := rtl.rt.RoundTrip(req)

	rtl.l.Infof("http: response: %s (request: %s %s)", res.Status, req.Method, req.URL.String())

	return res, err
}
