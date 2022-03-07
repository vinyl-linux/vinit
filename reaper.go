package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var (
	childCatcher = make(chan os.Signal, 1)
)

func init() {
	signal.Notify(childCatcher, syscall.SIGCHLD)
}

func reap() {
	for sig := range childCatcher {
		fmt.Printf("\n\n%#v\n\n", sig)
	}
}
