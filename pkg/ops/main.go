package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/soushin/bazel-multiprojects/pkg/ops/usecase"

	"github.com/soushin/bazel-multiprojects/pkg/ops/server"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/kelseyhightower/envconfig"
	"github.com/soushin/bazel-multiprojects/pkg/ops/client"
	"github.com/soushin/bazel-multiprojects/pkg/ops/handler"
	"github.com/soushin/bazel-multiprojects/proto/ops"
)

type envConfig struct {
	HttpPort string `envconfig:"HTTP_PORT" default:"8080"`
	GrpcPort string `envconfig:"GRPC_PORT" default:"50051"`

	GitHubToken string `envconfig:"GITHUB_TOKEN" default:"token"`
}

func main() {
	os.Exit(_main(os.Args[1:]))
}

func _main(args []string) int {
	appLog, _ := zap.NewProduction()
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		appLog.With(zap.Error(err)).Error("Failed to process env var")
		return 1
	}

	// client
	gitHubCli := client.NewGitHubClient(appLog, env.GitHubToken)

	// useCase
	deployUseCase := usecase.NewDeployUseCase(appLog, gitHubCli)

	// handler
	deployHandler := handler.NewDeployHandler(appLog, deployUseCase)

	// serve gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", env.GrpcPort))
	defer lis.Close()
	if err != nil {
		appLog.With(zap.Error(err)).Error("Failed to listen")
		return 1
	}

	grpcServer := grpc.NewServer()
	deployServer := server.NewDeployServer(appLog, deployHandler)
	ops.RegisterDeployServer(grpcServer, deployServer)
	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			appLog.With(zap.Error(err)).Error("Failed to serve gRPC Server")
		}
	}()
	appLog.With(zap.String("port", env.GrpcPort)).Info("Server listening")

	http.HandleFunc("/hc", hcHandler)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", env.HttpPort), nil); err != nil {
		appLog.With(zap.Error(err)).Error("Failed to serve http Server")
	}

	return 0
}

func hcHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "true")
}
