//go:build sudo
// +build sudo

package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestService_Run(t *testing.T) {
	d, _ := os.Getwd()

	s, err := LoadService("app-oneoff", filepath.Join(d, "testdata/services/00-app-oneoff"))
	if err != nil {
		t.Errorf("unexpected error %#v", err)
	}

	defer func() {
		err := recover()
		if err != nil {
			t.Errorf("unexpected error %#v", err)
		}
	}()

	err = s.Start(false)
	if err != nil {
		t.Errorf("unexpected error %#v", err)
	}

	time.Sleep(time.Millisecond * 200)

	if s.status.Error != nil {
		t.Errorf("unexpected error %#v", s.status.Error)
	}
}
