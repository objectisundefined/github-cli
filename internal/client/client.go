package client

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	config "github.com/objectisundefined/github-cli/internal/config"
)

// NewClient creates the Github client with token
func NewClient(ctx context.Context, cfg *config.Config) *github.Client {
	if len(cfg.Token) == 0 {
		return github.NewClient(nil)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.Token},
	)

	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}
