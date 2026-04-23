package github

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/google/go-github/v85/github"
	"github.com/suzuki-shunsuke/ghtkn-go-sdk/ghtkn"
	"golang.org/x/oauth2"
)

type Client struct {
	repo RepositoriesService
}

type RepositoriesService interface {
	GetContents(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (*github.RepositoryContent, []*github.RepositoryContent, *github.Response, error)
}

type (
	RepositoryContentGetOptions = github.RepositoryContentGetOptions
)

func New(ctx context.Context, logger *slog.Logger, token string, ghtknEnabled bool) *Client {
	gh := github.NewClient(getHTTPClient(ctx, logger, token, ghtknEnabled))
	return &Client{
		repo: gh.Repositories,
	}
}

func getHTTPClient(ctx context.Context, logger *slog.Logger, token string, ghtknEnabled bool) *http.Client {
	ts := getTokenSource(logger, token, ghtknEnabled)
	if ts == nil {
		return http.DefaultClient
	}
	return oauth2.NewClient(ctx, ts)
}

func getTokenSource(logger *slog.Logger, token string, ghtknEnabled bool) oauth2.TokenSource {
	if token != "" {
		return oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
	}
	if ghtknEnabled {
		return ghtkn.New().TokenSource(logger, &ghtkn.InputGet{})
	}
	return nil
}

func GetGitHubTokenFromEnv() string {
	for _, key := range []string{"DOCFRESH_GITHUB_TOKEN", "GITHUB_TOKEN"} {
		s := os.Getenv(key)
		if s != "" {
			return s
		}
	}
	return ""
}

func GetGHTKNEnabledFromEnv() (bool, error) {
	s := os.Getenv("DOCFRESH_GHTKN_ENABLED")
	if s == "" {
		return false, nil
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false, fmt.Errorf("parse the environment variable as a boolean: %w", err)
	}
	return b, nil
}
