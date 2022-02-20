package main

import (
	"context"

	"github.com/vinyl-linux/vinit/dispatcher"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	errNoService = status.Error(codes.InvalidArgument, "missing service name")
)

type Dispatcher struct {
	s *Supervisor

	dispatcher.UnimplementedDispatcherServer
}

func (d Dispatcher) Start(ctx context.Context, s *dispatcher.Service) (out *emptypb.Empty, err error) {
	out = new(emptypb.Empty)

	if s == nil || s.Name == "" {
		return out, errNoService
	}

	return out, d.s.Start(s.Name, false)
}

func (d Dispatcher) Stop(ctx context.Context, s *dispatcher.Service) (out *emptypb.Empty, err error) {
	out = new(emptypb.Empty)

	if s == nil || s.Name == "" {
		return out, errNoService
	}

	return out, d.s.Stop(s.Name)
}

func (d Dispatcher) Status(ctx context.Context, s *dispatcher.Service) (out *dispatcher.ServiceStatus, err error) {
	out = new(dispatcher.ServiceStatus)

	if s == nil || s.Name == "" {
		return out, errNoService
	}

	status, err := d.s.Status(s.Name)
	if err != nil {
		return
	}

	out.Svc = s
	out.Running = status.Running
	out.Pid = uint32(status.Pid)
	out.ExitStatus = uint32(status.ExitStatus)
	out.StartTime = timestamppb.New(status.StartTime)
	out.EndTime = timestamppb.New(status.EndTime)
	out.Success = status.Success

	if status.Error != nil {
		out.Error = status.Error.Error()
	}

	return
}

func (d Dispatcher) Reload(ctx context.Context, s *dispatcher.Service) (out *emptypb.Empty, err error) {
	out = new(emptypb.Empty)

	if s == nil || s.Name == "" {
		return out, errNoService
	}

	return out, d.s.Reload(s.Name)
}

func (d Dispatcher) ReadConfigs(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	return in, d.s.LoadConfigs()
}

func (d Dispatcher) SystemStatus(_ *emptypb.Empty, ds dispatcher.Dispatcher_SystemStatusServer) (err error) {
	var status *dispatcher.ServiceStatus

	for s := range d.s.services {
		status, err = d.Status(ds.Context(), &dispatcher.Service{Name: s})
		if err != nil {
			return
		}

		err = ds.Send(status)
		if err != nil {
			return
		}
	}

	return
}
