package main

import (
	"context"
	"testing"

	"google.golang.org/protobuf/types/known/emptypb"
)

type dummyRebooter struct {
	cmd int
}

func (d *dummyRebooter) Reboot(cmd int) error {
	d.cmd = cmd

	return nil
}

type dummySyncer struct {
	syncCount int
}

func (d *dummySyncer) Sync() {
	d.syncCount += 1
}

func TestDispatcher_Shutdown(t *testing.T) {
	d := newDispatcher()

	for _, test := range []struct {
		name       string
		cmd        func(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
		expectCmd  int
		expectSync int
	}{
		{"shutdown", d.Shutdown, poweroff, 1},
		{"reboot", d.Reboot, restart, 1},
		{"halt", d.Halt, halt, 0},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := new(dummyRebooter)
			rebooter = r.Reboot

			s := new(dummySyncer)
			syncer = s.Sync

			test.cmd(context.Background(), nil)

			if r.cmd != test.expectCmd {
				t.Errorf("expected %X, received %X", r.cmd, test.expectCmd)
			}

			if s.syncCount != test.expectSync {
				t.Errorf("expected %X, received %X", s.syncCount, test.expectSync)
			}
		})
	}

}
