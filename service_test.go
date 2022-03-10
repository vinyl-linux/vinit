package main

import (
	"testing"
	"time"
)

func TestLoadService(t *testing.T) {
	for _, test := range []struct {
		name        string
		dir         string
		expectError bool
	}{
		{"service dir does not exist", "testdata/nonsuch", true},
		{"env file is unusable", "testdata/erroring/env-is-dir", true},
		{"user does not exist", "testdata/erroring/nonesuch-user", true},
		{"user does not exist", "testdata/erroring/nonesuch-group", true},
		{"binary does not exist", "testdata/erroring/missing-bin", true},
		{"binary exists but is not a file", "testdata/erroring/nonregfile-bin", true},
		{"binary exists but is not executable", "testdata/erroring/nonexecfile-bin", true},
		{"happy path", "testdata/services/00-app", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := LoadService("", test.dir)

			if test.expectError && err == nil {
				t.Errorf("expected error, received none")
			} else if !test.expectError && err != nil {
				t.Errorf("unexpected error %#v", err)
			}
		})
	}
}

func TestService_Start_UnableToWriteLogs(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			t.Fatal(err)
		}
	}()

	s, err := LoadService("", "testdata/erroring/logs-unwritable")
	if err != nil {
		t.Fatal(err)
	}

	err = s.Start(false)
	if err != nil {
		t.Fatal(err)
	}

	// Give service time to start
	time.Sleep(time.Millisecond * 100)
}

func TestService_Start_UnableToCreateLogs(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			t.Fatal(err)
		}
	}()

	s, err := LoadService("", "testdata/erroring/logdir-uncreatable")
	if err != nil {
		t.Fatal(err)
	}

	err = s.Start(false)
	if err != nil {
		t.Fatal(err)
	}

	// Give service time to start
	time.Sleep(time.Millisecond * 100)
}
