package main

import (
	"context"
	"syscall"

	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	rebooter RebootFunc = syscall.Reboot
	syncer   SyncFunc   = syscall.Sync
)

const (
	// the following come from https://github.com/torvalds/linux/blob/master/include/uapi/linux/reboot.h
	restart  = 0x01234567
	poweroff = 0x4321FEDC
	halt     = 0xCDEF0123
)

// RebootFunc allows us to stub out syscall behaviour in tests, to avoid accidentally
// shutting down test machines
type RebootFunc func(int) error

// SyncFunc allows us to stub out calls to sync
type SyncFunc func()

// Shutdown will stop all services nicely, in reverse group/ priority order
// and then sends a shutdown signal to the kernel
func (d Dispatcher) Shutdown(context.Context, *emptypb.Empty) (out *emptypb.Empty, err error) {
	err = d.s.StopAll()
	if err != nil {
		return
	}

	syncer()

	return new(emptypb.Empty), rebooter(poweroff)
}

// Reboot will stop all services nicely, in reverse group/ priority order
// and then sends a reboots signal
func (d Dispatcher) Reboot(context.Context, *emptypb.Empty) (out *emptypb.Empty, err error) {
	err = d.s.StopAll()
	if err != nil {
		return
	}

	syncer()

	return new(emptypb.Empty), rebooter(restart)
}

// Halt will aggressively halt the system without bothering to stop anything or even
// syncing
func (d Dispatcher) Halt(context.Context, *emptypb.Empty) (out *emptypb.Empty, err error) {
	return new(emptypb.Empty), rebooter(halt)
}
