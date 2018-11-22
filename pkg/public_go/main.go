package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/soushin/bazel-multiprojects/pkg/common_go/util"
	"github.com/soushin/bazel-multiprojects/proto/echo"
	"github.com/soushin/bazel-multiprojects/proto/greet"
)

const (
	httpPort = "8080"
	grpcPort = "50051"
)

var msg string

type server struct{}

func (s *server) Greet(ctx context.Context, in *greet.GreetInbound) (*greet.GreetOutbound, error) {
	return &greet.GreetOutbound{
		Message: util.Add(fmt.Sprintf("%s Go!!", msg))}, nil
}

func (s *server) Echo(ctx context.Context, in *echo.EchoInbound) (*echo.EchoOutbound, error) {
	return &echo.EchoOutbound{
		Message: in.Message}, nil
}

func main() {
	g := flag.String("greet", "Hello", "greet message")
	flag.Parse()
	greetUsecase, err := initializeGreetUsecase(context.Background(), *g)
	if err != nil {
		log.Fatalln(err)
	}
	msg = greetUsecase.Msg

	// serve gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	defer lis.Close()
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	server := &server{}
	echo.RegisterEchoServer(grpcServer, server)
	greet.RegisterGreetServer(grpcServer, server)
	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve grpc server")
		}
	}()

	// serve http/1.1 server
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", httpPort), nil); err != nil {
		log.Fatalf("failed to serve http server")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, util.Add(fmt.Sprintf("%s Go!", msg)))
}
