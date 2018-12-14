package server

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/soushin/bazel-multiprojects/pkg/ops/handler"
	"github.com/soushin/bazel-multiprojects/proto/ops"
)

type DeployServer interface {
	GetTargets(ctx context.Context, in *empty.Empty) (*ops.TargetOutbound, error)
	Execute(inbound *ops.DeployInbound, stream ops.Deploy_ExecuteServer) error
}

type deployServerImpl struct {
	appLog  *zap.Logger
	handler handler.DeployHandler
}

func NewDeployServer(appLog *zap.Logger, handler handler.DeployHandler) DeployServer {
	return &deployServerImpl{
		appLog:  appLog,
		handler: handler,
	}
}

func (s *deployServerImpl) GetTargets(ctx context.Context, in *empty.Empty) (*ops.TargetOutbound, error) {

	owner := "soushin"
	repo := "bazel-multiprojects"
	packagePath := "pkg"

	paths, err := s.handler.Target(owner, repo, packagePath)
	if err != nil {
		s.appLog.With(zap.Strings("params", []string{owner, repo, packagePath})).Error("invalid process")
		return nil, errors.Wrap(err, "failed to get targets")
	}

	targets := make([]*ops.Target, len(paths))
	for i, p := range paths {
		targets[i] = &ops.Target{
			Repo: fmt.Sprintf("%s/%s", owner, repo),
			Pkg:  p,
		}
	}

	return &ops.TargetOutbound{
		Targets: targets,
	}, nil
}

func (s *deployServerImpl) Execute(inbound *ops.DeployInbound, stream ops.Deploy_ExecuteServer) error {

	owner := inbound.Owner
	repo := inbound.Repository
	packagePath := inbound.Package
	branch := inbound.Branch
	target := fmt.Sprintf("%s/%s:%s@%s", owner, repo, packagePath, branch)

	if err := stream.Send(&ops.DeployOutbound{
		Progress: ops.DeployProgress_STARTED,
		Message:  fmt.Sprintf("Started deploy: %s", target),
	}); err != nil {
		return err
	}

	if err := stream.Send(&ops.DeployOutbound{
		Progress: ops.DeployProgress_RUNNING,
		Message:  fmt.Sprintf("Running deploy: %s", target),
	}); err != nil {
		return err
	}

	err := s.handler.Execute(owner, repo, branch, packagePath)
	if err != nil {
		s.appLog.With(zap.Strings("params", []string{owner, repo, branch, packagePath})).Error("invalid process")
		if err := stream.Send(&ops.DeployOutbound{
			Progress: ops.DeployProgress_ERROR,
			Message:  err.Error(),
		}); err != nil {
			return err
		}
	}

	if err := stream.Send(&ops.DeployOutbound{
		Progress: ops.DeployProgress_SUCCESS,
		Message:  fmt.Sprintf("Succeess deploy: %s", target),
	}); err != nil {
		return err
	}

	return nil
}
