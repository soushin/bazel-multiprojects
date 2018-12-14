package handler

import (
	"go.uber.org/zap"

	"github.com/pkg/errors"
	"github.com/soushin/bazel-multiprojects/pkg/ops/usecase"
)

type DeployHandler interface {
	Target(owner, repo, path string) ([]string, error)
	Execute(owner, repo, branch, path string) error
}

type deployHandlerImpl struct {
	appLog  *zap.Logger
	useCase usecase.DeployUseCase
}

func NewDeployHandler(appLog *zap.Logger, useCase usecase.DeployUseCase) DeployHandler {
	return &deployHandlerImpl{
		appLog:  appLog,
		useCase: useCase,
	}
}

func (h *deployHandlerImpl) Target(owner, repo, path string) ([]string, error) {
	targets, err := h.useCase.GetContents(owner, repo, path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to useCase.GetContents")
	}
	return targets, nil
}

func (h *deployHandlerImpl) Execute(owner, repo, branch, path string) error {

	if err := h.useCase.ExistsContent(owner, repo, path); err != nil {
		return errors.Wrap(err, "failed to useCase.ExistsContent")
	}

	if err := h.useCase.ExistsBranch(owner, repo, branch); err != nil {
		return errors.Wrap(err, "failed to useCase.ExistsBranch")
	}

	checkoutPath, err := h.useCase.CheckoutBranch(owner, repo, branch)
	if err != nil {
		return errors.Wrap(err, "failed to useCase.CheckoutBranch")
	}

	if err := h.useCase.ReplaceImage(checkoutPath, path, owner, repo, branch); err != nil {
		return errors.Wrap(err, "failed to useCase.ReplaceImage")
	}

	return nil
}
