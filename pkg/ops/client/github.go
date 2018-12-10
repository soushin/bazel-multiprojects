package client

import (
	"context"

	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

type GitHubClient interface {
	GetContents(owner, repo, path string) ([]*github.RepositoryContent, error)
}

type gitHubClientImpl struct {
	appLog *zap.Logger
	token  string
}

func NewGitHubClient(appLog *zap.Logger, token string) GitHubClient {
	return &gitHubClientImpl{
		appLog: appLog,
		token:  token,
	}
}

func (c *gitHubClientImpl) GetContents(owner, repo, path string) ([]*github.RepositoryContent, error) {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.token},
	)
	tc := oauth2.NewClient(ctx, ts)

	cli := github.NewClient(tc)

	_, directoryContent, _, err := cli.Repositories.GetContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{})
	if err != nil {
		c.appLog.With(zap.String("owner", owner), zap.String("repo", repo), zap.String("path", path)).Error("invalid process")
		return nil, errors.Wrapf(err, "failed to get contents")
	}

	return directoryContent, nil
}
