package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	var err error

	maxLogLines = 10

	f, err := os.CreateTemp("", "")
	if err != nil {
		panic(err)
	}

	f.Close()

	sugar, err = NewLogger(f.Name())
	if err != nil {
		panic(err)
	}

	code := m.Run()
	os.Exit(code)
}

func TestSetup(t *testing.T) {
	svcDir = "testdata/services"

	_, err := Setup()
	if err != nil {
		t.Errorf("unexpected error %#v", err)
	}
}

func TestSetup_DodgyServicesJustWarns(t *testing.T) {
	svcDir = "testdata/mixed-status-services"

	_, err := Setup()
	if err != nil {
		t.Errorf("unexpected error %#v", err)
	}
}

func TestSetup_MissingSvcDir(t *testing.T) {
	svcDir = "/tmp/this/dir/hopefully/doesnt/exist"

	_, err := Setup()
	if err == nil {
		t.Errorf("expected error, received none")
	}
}

func TestSetup_MissingCerts(t *testing.T) {
	svcDir = "testdata/services"
	certDir = "/tmp/this/dir/hopefully/doesnt/exist"

	_, err := Setup()
	if err == nil {
		t.Errorf("expected error, received none")
	}
}
