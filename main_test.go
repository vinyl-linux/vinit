package main

import (
	"testing"
)

func TestSetup(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			t.Errorf("unexpected panic: %#v", err)
		}
	}()

	svcDir = "testdata/services"

	Setup()
}

func TestSetup_MissingSvcDir(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Error("expected panic, received none")
		}
	}()

	svcDir = "/tmp/this/dir/hopefully/doesnt/exist"

	Setup()
}

func TestSetup_MissingCerts(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Error("expected panic, received none")
		}
	}()

	svcDir = "testdata/services"
	certDir = "/tmp/this/dir/hopefully/doesnt/exist"

	Setup()
}
