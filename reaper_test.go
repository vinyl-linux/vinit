package main

import (
	"os"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

func TestReaper(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			t.Fatalf("unexpected panic %#v", err)
		}
	}()

	oldProcDir := procDir
	oldWaiterFunc := waiterFunc

	defer func() {
		procDir = oldProcDir
		waiterFunc = oldWaiterFunc
	}()

	procDir = "testdata/fake-proc"
	waiterFunc = func(int, *unix.WaitStatus, int, *unix.Rusage) (int, error) {
		return 0, nil
	}

	go reap()

	pid := os.Getpid()
	proc, err := os.FindProcess(pid)
	if err != nil {
		t.Fatalf(err.Error())
	}

	proc.Signal(unix.SIGCHLD)

	time.Sleep(time.Millisecond * 100)
}
