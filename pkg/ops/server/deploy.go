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
	path := "pkg"

	paths, err := s.handler.Target(owner, repo, path)
	if err != nil {
		s.appLog.With(zap.Strings("params", []string{owner, repo, path})).Error("invalid process")
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
