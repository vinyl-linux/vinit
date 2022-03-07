package main

import (
	"bufio"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

var (
	childCatcher        = make(chan os.Signal, 1)
	procDir             = "/proc"
	waiterFunc   waiter = unix.Wait4
)

type waiter func(int, *unix.WaitStatus, int, *unix.Rusage) (int, error)

func init() {
	signal.Notify(childCatcher, unix.SIGCHLD)
}

func reap() {
	var err error
	for range childCatcher {
		err = reapLoop()
		if err != nil {
			sugar.Errorw("reading proc failed",
				"error", err.Error(),
			)
		}

		// Bit of breathing room, maybe?
		time.Sleep(time.Millisecond * 100)
	}
}

func reapLoop() (err error) {
	pids, err := getZombiePids()
	if err != nil {
		return
	}

	if len(pids) == 0 {
		return
	}

	var (
		status = new(unix.WaitStatus)
		rusage = new(unix.Rusage)
	)

	sugar.Infow("reaper is reaping",
		"count", len(pids),
		"pids", pids,
	)

	for _, pid := range pids {
		_, err = waiterFunc(pid, status, unix.WNOHANG, rusage)
		if err != nil {
			break
		}
	}

	return
}

// getZombiePids runs through proc, looking for processes which are
// in a zombie state.
//
// It does this by finding any directory in `/proc/` whch starts with
// a number, and then passing that field to a function who can inspect
// that process (should it be a process) for zombie status, returning a
// slice containing any process in a zombie state.
func getZombiePids() (pids []int, err error) {
	f, err := os.Open(procDir)
	if err != nil {
		return
	}

	defer f.Close()

	pids = make([]int, 0)

	var (
		names []string
		isZ   bool
	)

	for {
		names, err = f.Readdirnames(20)
		if err != nil {
			if err == io.EOF {
				err = nil
			}

			break
		}

		for _, name := range names {
			if name[0] < '0' || name[0] > '9' {
				continue
			}

			isZ = checkZombie(name)
			if !isZ {
				continue
			}

			pid, err := strconv.ParseInt(name, 10, 0)
			if err != nil {
				continue
			}

			pids = append(pids, int(pid))
		}
	}

	return
}

// checkZombie reads the file /proc/${pid}/status, looking for whether
// or not the field 'State' signifies the process as a Zombie
//
// pid is a string here purely because at the time we check the pid, it
// makes sense to leave it as a string and not to strconv it until we need
// to.
func checkZombie(pid string) bool {
	f, err := os.Open(filepath.Join(procDir, pid, "status"))
	if err != nil {
		// swallow errors; maybe the pid has ended, or maybe there's another
		// issue that would stop us waiting on the process anyway /shrug
		return false
	}

	defer f.Close()

	var fields []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields = strings.Fields(scanner.Text())
		if fields[0] != "State:" {
			continue
		}

		if len(fields) > 1 && fields[1] == "Z" {
			return true
		}

		// At this point we've read the State line, and it wasn't
		// a zombie (or something weird happened).
		//
		// In anycase, sack it off
		break
	}

	return false
}
