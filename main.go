package main

import (
	"crypto/tls"
	"net"
	"os"
	"path/filepath"

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
	certDir  = "certs"
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

	srv, supervisor := Setup()

	if os.Getpid() == 1 {
		// try to delete sockAddr, if it exists
		os.Remove(sockAddr)
	}

	lis, err := net.Listen("unix", sockAddr)
	if err != nil {
		sugar.Panic(err)
	}

	go supervisor.RunShell()

	sugar.Panic(srv.Serve(lis))
}

func Setup() (grpcServer *grpc.Server, supervisor *Supervisor) {
	supervisor, err := New(svcDir)
	if err != nil {
		sugar.Panic(err)
	}

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		sugar.Panic(err)
	}

	supervisor.StartAll()

	d := Dispatcher{supervisor, dispatcher.UnimplementedDispatcherServer{}}

	grpcServer = grpc.NewServer(
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

	return
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
	serverCert, err := tls.LoadX509KeyPair(filepath.Join(certDir, "server-cert.pem"), filepath.Join(certDir, "server-key.pem"))
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
