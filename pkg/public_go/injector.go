//+build wireinject

package main

import (
	"context"

	"github.com/google/go-cloud/wire"

	"github.com/soushin/bazel-multiprojects/pkg/public_go/usecase"
)

//+build wireinject

func initializeGreetUsecase(ctx context.Context, greet string) (usecase.GreetUsecase, error) {
	wire.Build(usecase.GreetUsecaseSet)
	return usecase.GreetUsecase{}, nil
}
