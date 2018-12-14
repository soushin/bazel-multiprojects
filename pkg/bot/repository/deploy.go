package repository

import (
	"context"
	"io"

	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	"github.com/pkg/errors"
	"github.com/soushin/bazel-multiprojects/proto/ops"
)

type DeployRepository interface {
	Execute(owner, repo, path, branch string, observer func(result *ops.DeployOutbound)) error
}

type deployRepositoryImpl struct {
	appLog       *zap.Logger
	opsDeployCli ops.DeployClient
}

func NewDeployRepository(appLog *zap.Logger, opsDeployCli ops.DeployClient) DeployRepository {
	return &deployRepositoryImpl{
		appLog:       appLog,
		opsDeployCli: opsDeployCli,
	}
}

func (r *deployRepositoryImpl) Execute(owner, repo, path, branch string, observer func(result *ops.DeployOutbound)) error {

	md := metadata.New(map[string]string{})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	inbound := &ops.DeployInbound{
		Owner:      owner,
		Repository: repo,
		Package:    path,
		Branch:     branch,
	}

	stream, err := r.opsDeployCli.Execute(ctx, inbound)
	if err != nil {
		r.appLog.With(zap.Error(err)).Error("invalid process")
		return errors.Wrapf(err, "failed to deploy")
	}

	for {
		result, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		observer(result)
		if result.Progress == ops.DeployProgress_SUCCESS || result.Progress == ops.DeployProgress_ERROR {
			return stream.CloseSend()
		}
	}
}
