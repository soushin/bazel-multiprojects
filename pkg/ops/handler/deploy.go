package handler

import (
	"go.uber.org/zap"

	"github.com/pkg/errors"
	"github.com/soushin/bazel-multiprojects/pkg/ops/client"
)

type DeployHandler interface {
	Target(owner, repo, path string) ([]string, error)
}

type deployHandlerImpl struct {
	appLog    *zap.Logger
	githubCli client.GitHubClient
}

func NewDeployHandler(appLog *zap.Logger, githubCli client.GitHubClient) DeployHandler {
	return &deployHandlerImpl{
		appLog:    appLog,
		githubCli: githubCli,
	}
}

func (h *deployHandlerImpl) Target(owner, repo, path string) ([]string, error) {

	contents, err := h.githubCli.GetContents(owner, repo, path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to handle target")
	}

	targets := make([]string, len(contents))
	for i, content := range contents {
		targets[i] = *content.Path
	}

	return targets, nil
}
