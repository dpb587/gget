package gget

import (
	"net"
	"net/http"
	"os"
	"time"

	"github.com/dpb587/gget/pkg/app"
	"github.com/dpb587/gget/pkg/ggetutil"
	"github.com/sirupsen/logrus"
)

type Runtime struct {
	Help    bool                 `long:"help" short:"h" description:"show documentation of this command"`
	Quiet   bool                 `long:"quiet" short:"q" description:"suppress runtime status reporting"`
	Verbose []bool               `long:"verbose" short:"v" description:"increase logging verbosity (multiple)"`
	Version *ggetutil.VersionOpt `long:"version" description:"show version of this command (with optional constraint to validate)" optional:"true" optional-value:"*" value-name:"[CONSTRAINT]"`

	app        app.Version
	logger     *logrus.Logger
	httpClient *http.Client
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

func (r *Runtime) NewHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: roundTripLogger{
			l: r.Logger(),
			rt: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				Dial: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).Dial,
				IdleConnTimeout:       15 * time.Second,
				TLSHandshakeTimeout:   15 * time.Second,
				ResponseHeaderTimeout: 15 * time.Second,
				ExpectContinueTimeout: 5 * time.Second,
			},
		},
	}
}

type roundTripLogger struct {
	l  *logrus.Logger
	rt http.RoundTripper
}

func (rtl roundTripLogger) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	rtl.l.Debugf("http: %s %s", req.Method, req.URL.String())

	res, err := rtl.rt.RoundTrip(req)

	if res == nil {
		rtl.l.Infof("http: %s %s (response error)", req.Method, req.URL.String())
	} else {
		rtl.l.Infof("http: %s %s (status: %s)", req.Method, req.URL.String(), res.Status)

		if v := res.Header.Get("ratelimit-remaining"); v != "" {
			rtl.l.Debugf("http: %s %s (ratelimit-remaining: %s)", req.Method, req.URL.String(), v)
		} else if v := res.Header.Get("x-ratelimit-remaining"); v != "" {
			rtl.l.Debugf("http: %s %s (x-ratelimit-remaining: %s)", req.Method, req.URL.String(), v)
		}
	}

	return res, err
}
