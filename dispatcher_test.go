package main

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/vinyl-linux/vinit/dispatcher"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type dummyServiceStatusServer struct {
	grpc.ServerStream
	messages []*dispatcher.ServiceStatus
}

func (d *dummyServiceStatusServer) Send(m *dispatcher.ServiceStatus) error {
	d.messages = append(d.messages, m)

	return nil
}

func newDispatcher() Dispatcher {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	supervisor, err := New(pwd + "/testdata/services")
	if err != nil {
		panic(err)
	}

	return Dispatcher{supervisor, dispatcher.UnimplementedDispatcherServer{}}
}

func TestDispatcher_Stop(t *testing.T) {
	d := newDispatcher()

	defer func() {
		err := d.s.StopAll()
		if err != nil {
			t.Fatal(err)
		}
	}()

	t.Run("Happy path", func(t *testing.T) {
		_, err := d.Start(context.Background(), &dispatcher.Service{
			Name: "app",
		})
		if err != nil {
			t.Errorf("unexpected error: %#v", err)
		}

		time.Sleep(time.Millisecond * 100)

		_, err = d.Stop(context.Background(), &dispatcher.Service{
			Name: "app",
		})
		if err != nil {
			t.Errorf("unexpected error: %#v", err)
		}

		time.Sleep(time.Millisecond * 200)

		if d.s.services["app"].isRunning() {
			t.Error("app either did not stop, or did not update status")
		}
	})

	t.Run("Service must not be empty", func(t *testing.T) {
		_, err := d.Stop(context.Background(), &dispatcher.Service{
			Name: "",
		})
		if err == nil || err.Error() != "rpc error: code = InvalidArgument desc = missing service name" {
			t.Errorf("error: %#v should be %q", err.Error(), "missing service name")
		}
	})

	t.Run("Service must exist", func(t *testing.T) {
		_, err := d.Stop(context.Background(), &dispatcher.Service{
			Name: "foo",
		})
		if err == nil || err.Error() != "rpc error: code = InvalidArgument desc = service does not exist" {
			t.Errorf("error: %#v should be %q", err, "service does not exist")
		}
	})

	t.Run("Cannot stop a stopped service", func(t *testing.T) {
		_, err := d.Stop(context.Background(), &dispatcher.Service{
			Name: "app",
		})
		if err == nil || err.Error() != "service is not running" {
			t.Errorf("error: %#v should be %q", err, "service is not running")
		}
	})
}

func TestDispatcher_Start(t *testing.T) {
	d := newDispatcher()

	defer func() {
		time.Sleep(time.Millisecond * 100)

		err := d.s.StopAll()
		if err != nil {
			t.Fatal(err)
		}
	}()

	t.Run("Happy path", func(t *testing.T) {
		_, err := d.Start(context.Background(), &dispatcher.Service{
			Name: "app",
		})
		if err != nil {
			t.Errorf("unexpected error: %#v", err)
		}
	})

	// Give time for 'app' service to start and for supervisor to update
	// status
	time.Sleep(time.Millisecond * 100)

	t.Run("Cannot start service more than once", func(t *testing.T) {
		_, err := d.Start(context.Background(), &dispatcher.Service{
			Name: "app",
		})
		if err == nil || err.Error() != "service is already running" {
			t.Errorf("error: %#v should be %q", err, "service is already running")
		}
	})

	t.Run("Service must not be empty", func(t *testing.T) {
		_, err := d.Start(context.Background(), &dispatcher.Service{
			Name: "",
		})
		if err == nil || err.Error() != "rpc error: code = InvalidArgument desc = missing service name" {
			t.Errorf("error: %#v should be %q", err.Error(), "missing service name")
		}
	})

	t.Run("Service must exist", func(t *testing.T) {
		_, err := d.Start(context.Background(), &dispatcher.Service{
			Name: "foo",
		})
		if err == nil || err.Error() != "rpc error: code = InvalidArgument desc = service does not exist" {
			t.Errorf("error: %#v should be %q", err, "service does not exist")
		}
	})
}

func TestDispatcher_Status(t *testing.T) {
	d := newDispatcher()

	defer func() {
		time.Sleep(time.Millisecond * 100)

		err := d.s.StopAll()
		if err != nil {
			t.Fatal(err)
		}
	}()

	t.Run("Service is not running", func(t *testing.T) {
		s, err := d.Status(context.Background(), &dispatcher.Service{
			Name: "app",
		})
		if err != nil {
			t.Errorf("unexpected error: %#v", err)
		}

		if s.Running != false || s.StartTime.String() != s.EndTime.String() {
			t.Errorf("status does not appear to be empty %#v", s)
		}
	})

	t.Run("Service is running", func(t *testing.T) {
		_, err := d.Start(context.Background(), &dispatcher.Service{
			Name: "app",
		})
		if err != nil {
			t.Errorf("unexpected error: %#v", err)
		}

		time.Sleep(time.Millisecond * 100)

		s, err := d.Status(context.Background(), &dispatcher.Service{
			Name: "app",
		})
		if err != nil {
			t.Errorf("unexpected error: %#v", err)
		}

		if s.Running == false || s.StartTime == s.EndTime {
			t.Errorf("status appears to be empty %#v", s)
		}

		_, err = d.Stop(context.Background(), &dispatcher.Service{
			Name: "app",
		})
		if err != nil {
			t.Errorf("unexpected error: %#v", err)
		}

		time.Sleep(time.Millisecond * 100)
	})

	t.Run("Service is stopped", func(t *testing.T) {
		_, err := d.Start(context.Background(), &dispatcher.Service{
			Name: "app",
		})
		if err != nil {
			t.Errorf("unexpected error: %#v", err)
		}

		time.Sleep(time.Millisecond * 100)

		_, err = d.Stop(context.Background(), &dispatcher.Service{
			Name: "app",
		})
		if err != nil {
			t.Errorf("unexpected error: %#v", err)
		}

		time.Sleep(time.Millisecond * 200)

		s, err := d.Status(context.Background(), &dispatcher.Service{
			Name: "app",
		})
		if err != nil {
			t.Errorf("unexpected error: %#v", err)
		}

		if s.Running != false {
			t.Errorf("status contains old state %#v", s)
		}
	})

	t.Run("Service must not be empty", func(t *testing.T) {
		_, err := d.Status(context.Background(), &dispatcher.Service{
			Name: "",
		})
		if err == nil || err.Error() != "rpc error: code = InvalidArgument desc = missing service name" {
			t.Errorf("error: %#v should be %q", err.Error(), "missing service name")
		}
	})

	t.Run("Service must exist", func(t *testing.T) {
		_, err := d.Status(context.Background(), &dispatcher.Service{
			Name: "foo",
		})
		if err == nil || err.Error() != "rpc error: code = InvalidArgument desc = service does not exist" {
			t.Errorf("error: %#v should be %q", err, "service does not exist")
		}
	})

}

func TestDispatcher_Reload(t *testing.T) {
	d := newDispatcher()

	defer func() {
		time.Sleep(time.Millisecond * 100)

		err := d.s.StopAll()
		if err != nil {
			t.Fatal(err)
		}
	}()

	t.Run("Happy path", func(t *testing.T) {
		_, err := d.Start(context.Background(), &dispatcher.Service{
			Name: "app",
		})
		if err != nil {
			t.Errorf("unexpected error: %#v", err)
		}

		time.Sleep(time.Millisecond * 100)

		_, err = d.Reload(context.Background(), &dispatcher.Service{
			Name: "app",
		})
		if err != nil {
			t.Errorf("unexpected error: %#v", err)
		}

		_, err = d.Stop(context.Background(), &dispatcher.Service{
			Name: "app",
		})
		if err != nil {
			t.Errorf("unexpected error: %#v", err)
		}

		time.Sleep(time.Millisecond * 100)
	})

	t.Run("Cannot reload non-started service", func(t *testing.T) {
		_, err := d.Reload(context.Background(), &dispatcher.Service{
			Name: "app",
		})
		if err == nil || err.Error() != "service is not running" {
			t.Errorf("error: %#v should be %q", err, "service is not running")
		}
	})

	t.Run("Service must not be empty", func(t *testing.T) {
		_, err := d.Reload(context.Background(), &dispatcher.Service{
			Name: "",
		})
		if err == nil || err.Error() != "rpc error: code = InvalidArgument desc = missing service name" {
			t.Errorf("error: %#v should be %q", err.Error(), "missing service name")
		}
	})

	t.Run("Service must exist", func(t *testing.T) {
		_, err := d.Reload(context.Background(), &dispatcher.Service{
			Name: "foo",
		})
		if err == nil || err.Error() != "rpc error: code = InvalidArgument desc = service does not exist" {
			t.Errorf("error: %#v should be %q", err, "service does not exist")
		}
	})

}

func TestDispatcher_ReadConfigs(t *testing.T) {
	d := newDispatcher()

	defer func() {
		time.Sleep(time.Millisecond * 100)

		err := d.s.StopAll()
		if err != nil {
			t.Fatal(err)
		}
	}()

	_, err := d.Start(context.Background(), &dispatcher.Service{
		Name: "app",
	})
	if err != nil {
		t.Errorf("unexpected error: %#v", err)
	}

	time.Sleep(time.Millisecond * 100)

	currentStatus := d.s.services["app"].status

	_, err = d.ReadConfigs(context.Background(), new(emptypb.Empty))
	if err != nil {
		t.Errorf("unexpected error: %#v", err)
	}

	if !reflect.DeepEqual(currentStatus, d.s.services["app"].status) {
		t.Errorf("expected %#v, received %#v", currentStatus, d.s.services["app"].status)
	}
}

func TestDispatcher_Version(t *testing.T) {
	d := newDispatcher()

	vm, err := d.Version(context.Background(), nil)
	if err != nil {
		t.Errorf("unexpected error: %#v", err)
	}

	expect := &dispatcher.VersionMessage{
		Ref:       ref,
		BuildUser: buildUser,
		BuiltOn:   builtOn,
	}

	if !reflect.DeepEqual(expect, vm) {
		t.Errorf("expected %#v, received %#v", expect, vm)
	}
}

func TestDispatcher_SystemStatus(t *testing.T) {
	d := newDispatcher()

	dss := new(dummyServiceStatusServer)

	err := d.SystemStatus(new(emptypb.Empty), dss)
	if err != nil {
		t.Errorf("unexpected error: %#v", err)
	}

	if len(d.s.services) != len(dss.messages) {
		t.Errorf("expected %d messages, received %d", len(d.s.services), len(dss.messages))
	}
}
