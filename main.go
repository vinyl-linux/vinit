package main

import (
	"crypto/tls"
	"net"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/vinyl-linux/vinit/dispatcher"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger

	sockAddr = "/run/vinit.sock"
	svcDir   = envOrDefault("SVC_DIR", "/etc/vinit/services")
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

	srv := Setup()

	lis, err := net.Listen("tcp", sockAddr)
	if err != nil {
		sugar.Panic(err)
	}

	sugar.Panic(srv.Serve(lis))
}

func Setup() *grpc.Server {
	supervisor, err := New(svcDir)
	if err != nil {
		panic(err)
	}

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		sugar.Panic(err)
	}

	err = supervisor.StartAll()
	if err != nil {
		panic(err)
	}

	d := Dispatcher{supervisor, dispatcher.UnimplementedDispatcherServer{}}

	grpcServer := grpc.NewServer(
		grpc.Creds(tlsCredentials),
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

func envOrDefault(envvar, def string) string {
	out, ok := os.LookupEnv(envvar)
	if ok {
		return out
	}

	return def
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("certs/server-cert.pem", "certs/server-key.pem")
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
		MinVersion:   tls.VersionTLS12,
	}

	return credentials.NewTLS(config), nil
}
