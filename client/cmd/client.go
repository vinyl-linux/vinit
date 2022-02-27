package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/url"
	"path/filepath"
	"time"

	vinit "github.com/vinyl-linux/vinit/dispatcher"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/emptypb"
)

type client struct {
	c vinit.DispatcherClient
}

func newClient(addr string) (c client, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resolvedAddr, err := parseAddr(addr)
	if err != nil {
		return
	}

	conn, err := grpc.DialContext(ctx,
		resolvedAddr,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true,
		})),
	)

	if err != nil {
		return
	}

	c.c = vinit.NewDispatcherClient(conn)

	return
}

func parseAddr(addr string) (s string, err error) {
	u, err := url.Parse(addr)
	if err != nil {
		return
	}

	if u.Scheme == "unix" {
		p := u.Path

		u.Path, err = filepath.EvalSymlinks(p)
		if err != nil {
			return
		}
	}

	return u.String(), nil
}

func (c client) start(svc string) (err error) {
	is := &vinit.Service{
		Name: svc,
	}

	_, err = c.c.Start(context.Background(), is)

	return
}

func (c client) stop(svc string) (err error) {
	is := &vinit.Service{
		Name: svc,
	}

	_, err = c.c.Stop(context.Background(), is)

	return
}

func (c client) status(svc string) (status *vinit.ServiceStatus, err error) {
	is := &vinit.Service{
		Name: svc,
	}

	return c.c.Status(context.Background(), is)
}

func (c client) systemStatus() (statuses []*vinit.ServiceStatus, err error) {
	statuses = make([]*vinit.ServiceStatus, 0)

	sc, err := c.c.SystemStatus(context.Background(), new(emptypb.Empty))
	if err != nil {
		return
	}

	var status *vinit.ServiceStatus

	for {
		status, err = sc.Recv()
		if err != nil {
			if err == io.EOF {
				err = nil
			}

			break
		}

		statuses = append(statuses, status)
	}

	return
}

func (c client) version() (out string, err error) {
	v, err := c.c.Version(context.Background(), new(emptypb.Empty))
	if err != nil {
		return
	}

	return formatVersion(true, v.Ref, v.BuildUser, v.BuiltOn), nil
}

func (c client) reboot() (err error) {
	_, err = c.c.Reboot(context.Background(), new(emptypb.Empty))

	return
}

func formatVersion(isServer bool, ref, user, built string) string {
	return fmt.Sprintf("%s version\n---\nVersion: %s\nBuild User: %s\nBuilt On: %s\n",
		isServerString(isServer), ref, user, built,
	)
}

func isServerString(b bool) string {
	if b {
		return "Server"
	}

	return "Client"
}
