package ghet

import (
	"context"
	"net/http"

	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
)

type Global struct {
	Quiet  bool   `long:"quiet" description:"suppress status reporting"`
	Server string `long:"server" description:"use a custom GitHub Server" env:"GITHUB_SERVER"`
	Token  string `long:"token" description:"use a specific GitHub authentication token" env:"GITHUB_TOKEN"`
}

func (g Global) GitHubClient(server string) *github.Client {
	var tc *http.Client
	ctx := context.Background()

	if v := g.Token; v != "" {
		tc = oauth2.NewClient(
			ctx,
			oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: v},
			),
		)
	}

	return github.NewClient(tc)
}
