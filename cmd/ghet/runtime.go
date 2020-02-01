package ghet

import (
	"context"
	"net/http"
	"os"

	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
)

type Runtime struct {
	Quiet    bool `long:"quiet" description:"suppress status reporting"`
	Parallel int  `long:"parallel" description:"maximum number of parallel operations" default:"3"`

	Help    bool `long:"help" short:"h" description:"show documentation of this tool"`
	Version bool `long:"version" description:"show version of this tool"`
}

func (r Runtime) GitHubClient(server string) *github.Client {
	var tc *http.Client
	ctx := context.Background()

	if v := os.Getenv("GITHUB_TOKEN"); v != "" {
		tc = oauth2.NewClient(
			ctx,
			oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: v},
			),
		)
	}

	return github.NewClient(tc)
}
