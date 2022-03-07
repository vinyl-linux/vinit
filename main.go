package main

import (
	"crypto/tls"
	"net"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/vinyl-linux/vinit/dispatcher"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

var (
	sockAddr = "/run/vinit.sock"
	svcDir   = envOrDefault("SVC_DIR", "/etc/vinit/services")
	certDir  = "certs"
)

func main() {
	go reap()

	defer os.Remove(sockAddr)

	srv, err := Setup()
	if err != nil {
		sugar.Errorw("setup failed, booting into recovery shell",
			"error", err.Error(),
		)

		recoveryShell()

		return
	}

	if os.Getpid() == 1 {
		// try to delete sockAddr, if it exists.
		//
		// we don't care about the outcome of this; if the file doesn't
		// exist then happy days, otherwise we'll get a more useful error
		// when we try to listen anyway
		os.Remove(sockAddr) //#nosec: G104
	}

	lis, err := net.Listen("unix", sockAddr)
	if err != nil {
		sugar.Errorw("could not listen on socket address, booting into recovery shell",
			"sockAddr", sockAddr,
			"error", err.Error(),
		)

		recoveryShell()

		return
	}

	err = srv.Serve(lis)
	sugar.Errorw("vinit failed",
		"error", err.Error(),
	)

	recoveryShell()
}

func Setup() (grpcServer *grpc.Server, err error) {
	var supervisor *Supervisor

	supervisor, err = New(svcDir)
	if err != nil {
		if _, ok := err.(ConfigParseError); !ok {
			return
		}

		sugar.Warnw("could not load all configs",
			"error", err.Error(),
		)

		err = nil
	}

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		return
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

func recoveryShell() {
	logger.Info("Press Ctrl+D to reboot")

	c := exec.Command("/sbin/agetty", "-L", "-8", "--autologin", "root", "115200", "tty1", "linux")
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	_ = c.Run()
}
