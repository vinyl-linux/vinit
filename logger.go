package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/sys/unix"
)

var (
	sugar Logger

	maxLogLines = 1024
)

type Logger struct {
	Buffer []string

	c chan string
	f io.ReadWriter
}

func NewLogger(kmesgF string) (l Logger, err error) {
	l.Buffer = make([]string, maxLogLines)

	l.c = make(chan string)
	l.f, err = os.OpenFile(kmesgF, os.O_RDWR|unix.O_CLOEXEC|unix.O_NONBLOCK|unix.O_NOCTTY, 0o666) // #nosec: G302,G304
	if err != nil {
		return
	}

	go l.Start()

	return
}

func (l *Logger) Start() {
	for msg := range l.c {
		if l.f != nil {
			fmt.Fprint(l.f, msg)
		}

		l.Buffer = append([]string{msg}, l.Buffer[:maxLogLines-1]...)
	}
}

func (l Logger) Infow(msg string, kvs ...interface{}) {
	l.addLog("info", msg, kvs...)
}

func (l Logger) Errorw(msg string, kvs ...interface{}) {
	l.addLog("error", msg, kvs...)
}

func (l Logger) Warnw(msg string, kvs ...interface{}) {
	l.addLog("warn", msg, kvs...)
}

func (l Logger) addLog(level, msg string, kvs ...interface{}) {
	elems := make([]string, 0)
	elems = append(elems, fmt.Sprintf("vinit %s: %s", level, msg))

	switch len(kvs) {
	case 0, 2:
	case 1:
		kvs = nil
	default:
		if len(kvs)%2 != 0 {
			kvs = kvs[:len(kvs)-1]
		}
	}

	for i := 0; i < len(kvs); i += 2 {
		elems = append(elems, fmt.Sprintf(`%v="%v"`, kvs[i], kvs[i+1]))
	}

	l.c <- strings.Join(elems, " ")
}
