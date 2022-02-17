package main

import (
	"net"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/vinyl-linux/vinit/dispatcher"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger

	sockAddr = "/tmp/vinit.sock"
	svcDir   = os.Getenv("SVC_DIR")
)

func init() {
	var err error

	if logger == nil {
		logger, err = zap.NewProduction()
		if err != nil {
			panic(err)
		}
	}

	sugar = logger.Sugar()
}

func main() {
	defer os.Remove(sockAddr)

	lis, err := net.Listen("unix", sockAddr)
	if err != nil {
		sugar.Panic(err)
	}

	sugar.Panic(Setup().Serve(lis))
}

func Setup() *grpc.Server {
	supervisor, err := New(svcDir)
	if err != nil {
		panic(err)
	}

	d := Dispatcher{supervisor, dispatcher.UnimplementedDispatcherServer{}}

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(logger),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(logger),
		)),
	)

	dispatcher.RegisterDispatcherServer(grpcServer, d)
	reflection.Register(grpcServer)

	return grpcServer
}
