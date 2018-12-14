package handler

import (
	"go.uber.org/zap"

	"github.com/pkg/errors"
	"github.com/soushin/bazel-multiprojects/pkg/ops/usecase"
)

type DeployHandler interface {
	Target(owner, repo, packagePath string) ([]string, error)
	Execute(owner, repo, branch, packagePath string) error
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

func (h *deployHandlerImpl) Target(owner, repo, packagePath string) ([]string, error) {
	targets, err := h.useCase.GetContents(owner, repo, packagePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to useCase.GetContents")
	}
	return targets, nil
}

func (h *deployHandlerImpl) Execute(owner, repo, branch, packagePath string) error {

	if err := h.useCase.ExistsContent(owner, repo, packagePath); err != nil {
		return errors.Wrap(err, "failed to useCase.ExistsContent")
	}

	if err := h.useCase.ExistsBranch(owner, repo, branch); err != nil {
		return errors.Wrap(err, "failed to useCase.ExistsBranch")
	}

	checkoutPath, err := h.useCase.CheckoutBranch(owner, repo, branch)
	if err != nil {
		return errors.Wrap(err, "failed to useCase.CheckoutBranch")
	}

	if err := h.useCase.ReplaceImage(checkoutPath, packagePath, owner, repo, branch); err != nil {
		return errors.Wrap(err, "failed to useCase.ReplaceImage")
	}

	if err := h.useCase.Build(checkoutPath, packagePath); err != nil {
		return errors.Wrap(err, "failed to useCase.Build")
	}

	return nil
}
